package com.testprobe.app;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;

import com.sun.net.httpserver.HttpContext;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;

public class App {
  public static void main(String[] args) throws IOException {
      System.out.println("Starting simple java http server");
      HttpServer server = HttpServer.create(new InetSocketAddress(28500), 0);
      HttpContext context = server.createContext("/");
      context.setHandler(App::handleRequest);
      server.start();
  }

  private static void handleRequest(HttpExchange exchange) throws IOException {
      String response = "Hi there!";
      exchange.sendResponseHeaders(200, response.getBytes().length);//response code and length
      OutputStream os = exchange.getResponseBody();
      os.write(response.getBytes());
      os.close();
  }
}
