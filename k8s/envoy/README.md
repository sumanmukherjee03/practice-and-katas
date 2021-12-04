The envoy binary needs to be compiled. So, instead we would be using an official docker image provided by envoy for running the binary.

```
docker image build -t envoy-demo .
docker container run --name envoy-demo-test -p 10000:10000 -p 9901:9901 -d envoy-demo
```

Now, if you visit http://localhost:10000/google or http://localhost:10000/vpp
you will see the reverse proxy working.
There will be some 404s for urls in those pages, like for css and images and js hosted by google or virtualpairprogrammers.
But those 404s are coming from envoy and not from the corresponding targets.
That's why in the envoy admin interface, if you look at the prometheus metrics you would not
find the metrics for 404. You will only find metrics for 2XX.

Now, if you want to see some 404 prometheus metrics you can always try hitting http://localhost:10000/google/unknown
or some other non-existent url like that which will get forwarded to google.

To stop and remove the test container and the image
```
docker stop envoy-demo-test
docker rm envoy-demo-test
docker rmi envoy-demo
```

With istio, all the features of envoy are supported but we dont have to learn about envoy in depth.
The istiod pod takes care of injecting envoy config based upon the CRDs we configure in yaml files for istio.
More specifically this is the work of the "pilot" component in istiod.
