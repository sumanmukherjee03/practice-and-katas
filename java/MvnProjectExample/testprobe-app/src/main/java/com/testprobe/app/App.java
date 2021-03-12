package com.testprobe.app;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.util.*;
import java.util.concurrent.*;
import javassist.*;

import com.sun.net.httpserver.HttpContext;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;

public class App {
  private static Random rand = new Random();
  private static javassist.ClassPool cp = javassist.ClassPool.getDefault();
  private static final float CAP = 0.8f;  // 80%
  private static final int ONE_MB = 1024 * 1024;
  private static final Vector cache = new Vector();
  private static final Runtime rt = Runtime.getRuntime();
  private static final ArrayList<Double> list = new ArrayList<Double>(1000000);

  public static void main(String[] args) throws IOException {
    System.out.println("Starting simple java http server");

    long maxMemBytes = rt.maxMemory();
    long usedMemBytes = rt.totalMemory() - rt.freeMemory();
    long freeMemBytes = rt.maxMemory() - usedMemBytes;
    int allocBytes = Math.round(freeMemBytes * CAP);
    System.out.println("Initial max memory : " + maxMemBytes/ONE_MB + "MB");
    System.out.println("Initial used memory: " + usedMemBytes/ONE_MB + "MB");
    System.out.println("Initial free memory: " + freeMemBytes/ONE_MB + "MB");
    System.out.println("Reserve: " + allocBytes/ONE_MB + "MB");

    final Executor multi = Executors.newCachedThreadPool();
    HttpServer server = HttpServer.create(new InetSocketAddress(28500), 20);
    server.setExecutor(multi);

    HttpContext healthzContext = server.createContext("/healthz");
    healthzContext.setHandler(App::healthzHandler);

    HttpContext containerOOMContext = server.createContext("/containerOOMContext");
    containerOOMContext.setHandler(App::containerOOMHandler);

    HttpContext hsOOMContext = server.createContext("/hsOOMContext");
    hsOOMContext.setHandler(App::hsOOMHandler);

    HttpContext gcOverheadContext = server.createContext("/gcOverheadContext");
    gcOverheadContext.setHandler(App::gcOverheadHandler);

    HttpContext metaspaceOOMContext = server.createContext("/metaspaceOOMContext");
    metaspaceOOMContext.setHandler(App::metaspaceOOMHandler);

    HttpContext infoContext = server.createContext("/statusInfo");
    infoContext.setHandler(App::statusInfoHandler);

    server.start();
  }

  // Handler that when triggered repeatedly can lead to GC overhead, eventually causing a java.lang.OutOfMemoryError
  // The GC stays busy all the time trying to clean up heapspace while the application itself needs heapspace to allocate memory constantly to function properly.
  // This constant race between cleaning up heapspace and allocating more memory from the heap will lead to GC overhead error.
  // Controlled by increasing the heapspace with -Xmx512m parameter and -XX:-UseGCOverheadLimit parameter.
  private static void statusInfoHandler(HttpExchange exchange) throws IOException {
    System.out.println("status info called");
    long maxMemBytes = rt.maxMemory();
    long usedMemBytes = rt.totalMemory() - rt.freeMemory();
    long freeMemBytes = rt.maxMemory() - usedMemBytes;
    int allocBytes = Math.round(freeMemBytes * CAP);
    System.out.println("Current max memory : " + maxMemBytes/ONE_MB + "MB");
    System.out.println("Current used memory: " + usedMemBytes/ONE_MB + "MB");
    System.out.println("Current free memory: " + freeMemBytes/ONE_MB + "MB");
    System.out.println("Current Reserve: " + allocBytes/ONE_MB + "MB");
    String response = "SUCCESS";
    exchange.sendResponseHeaders(200, response.getBytes().length); //response code and length
    OutputStream os = exchange.getResponseBody();
    os.write(response.getBytes());
    os.close();
  }

  // Healthcheck that a docker/kubernetes/ecs etc can hit for keeping the service alive
  // To simulate a memory leak that manifests itself over time, this healthcheck should be sufficient.
  // It keeps filling up the final list object a little bit at a time and this object is not GC'd.
  private static void healthzHandler(HttpExchange exchange) throws IOException {
    System.out.println("Healthcheck called");
    for (int i = 0; i < 1000; i++){
      list.add(rand.nextDouble()); // make the healthcheck allocate some space because that will keep the heap space slowly filling up
    }
    String response = "SUCCESS";
    exchange.sendResponseHeaders(200, response.getBytes().length);//response code and length
    OutputStream os = exchange.getResponseBody();
    os.write(response.getBytes());
    os.close();
  }

  // Handler that when triggered repeatedly can lead to heapspace getting filled up causing an java.lang.OutOfMemoryError
  // In this method we are trying to allocate a big object and if there is not enough space in the heap it will trigger
  // an out of memory error and cause a heapdump if enabled
  private static void hsOOMHandler(HttpExchange exchange) throws IOException {
    System.out.println("hsOOMContext called");
    long[][] ary = new long[Integer.MAX_VALUE][Integer.MAX_VALUE];
    // Integer[] array = new Integer[20000 * 20000];
    String response = "FAILURE";
    exchange.sendResponseHeaders(500, response.getBytes().length);//response code and length
    OutputStream os = exchange.getResponseBody();
    os.write(response.getBytes());
    os.close();
  }

  // Handler that when triggered repeatedly can lead to heapspace getting filled up 25 MB at a time
  // Since this is filling up a Vector which is declared as final, it wont be GC'd
  // Over time if you trigger it multiple times it will fill up the heap upto the max allocated heapspace
  // Any subsequent operation that requires allocating memory will causing a java.lang.OutOfMemoryError
  private static void containerOOMHandler(HttpExchange exchange) throws IOException {
    System.out.println("containerOOMContext called");
    for (int i = 0; i < 25; i++){
      cache.add(new byte[ONE_MB]);
    }
    long maxMemBytes = rt.maxMemory();
    long usedMemBytes = rt.totalMemory() - rt.freeMemory();
    long freeMemBytes = rt.maxMemory() - usedMemBytes;
    int allocBytes = Math.round(freeMemBytes * CAP);
    System.out.println("Current max memory : " + maxMemBytes/ONE_MB + "MB");
    System.out.println("Current used memory: " + usedMemBytes/ONE_MB + "MB");
    System.out.println("Current free memory: " + freeMemBytes/ONE_MB + "MB");
    System.out.println("Current Reserve: " + allocBytes/ONE_MB + "MB");
    String response = "SUCCESS";
    exchange.sendResponseHeaders(200, response.getBytes().length); //response code and length
    OutputStream os = exchange.getResponseBody();
    os.write(response.getBytes());
    os.close();
  }

  // Handler that when triggered repeatedly can lead to GC overhead, eventually causing a java.lang.OutOfMemoryError
  // The GC stays busy all the time trying to clean up heapspace while the application itself needs heapspace to allocate memory constantly to function properly.
  // This constant race between cleaning up heapspace and allocating more memory from the heap will lead to GC overhead error.
  // Controlled by increasing the heapspace with -Xmx512m parameter and -XX:-UseGCOverheadLimit parameter.
  private static void gcOverheadHandler(HttpExchange exchange) throws IOException {
    System.out.println("gcOverheadContext called");
    Map m = new HashMap();
    m = System.getProperties();
    while (true) {
      m.put(rand.nextInt(), "randomValue");
    }
  }

  // Handler that when triggered can cause the metaspace and permgen space to fillup quickly causing a java.lang.OutOfMemoryError
  // Metaspace is allocated from the same address space as the heapspace, so if you increase this, then it simultaneously reduces the heapspace.
  // Controlled by -XX:MaxMetaSpaceSize=512m parameter and -XX:MaxPermSize=512m parameter
  // We are filling up the metaspace here by loading more and more classes.
  // If there is no max meta space size set, the metaspace will continue to grow until it hits the limit of the container allocated memory.
  // In such a case you will see a container OOM error but no java OOM error in the logs.
  private static void metaspaceOOMHandler(HttpExchange exchange) throws IOException {
    System.out.println("metaspaceOOMContext called");
    for (int i = 0; i < 100000; i++) {
      javassist.CtClass cc = cp.makeClass("com.testprobe.app.Metaspace" + i);
      try {
        Class c = cc.toClass();
      } catch(Exception e) {
        throw new IOException("Class could not be created");
      }
    }
    long maxMemBytes = rt.maxMemory();
    long usedMemBytes = rt.totalMemory() - rt.freeMemory();
    long freeMemBytes = rt.maxMemory() - usedMemBytes;
    int allocBytes = Math.round(freeMemBytes * CAP);
    System.out.println("Current max memory : " + maxMemBytes/ONE_MB + "MB");
    System.out.println("Current used memory: " + usedMemBytes/ONE_MB + "MB");
    System.out.println("Current free memory: " + freeMemBytes/ONE_MB + "MB");
    System.out.println("Current Reserve: " + allocBytes/ONE_MB + "MB");
    String response = "SUCCESS";
    exchange.sendResponseHeaders(200, response.getBytes().length);//response code and length
    OutputStream os = exchange.getResponseBody();
    os.write(response.getBytes());
    os.close();
  }
}
