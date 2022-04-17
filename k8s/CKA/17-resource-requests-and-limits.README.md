## resource requests and limits
Here's another example of a pod-definition.yaml but with resource requests
0.1 cpu == 100m where m stands for milli.
cpu = 1 means 1 vCPU in AWS or GCP or 1 Hyperthread.
mem = 256Mi where Mi stands for Mebibyte.
For memory G is gigabyte, Gi is gibibyte.

By default kubernetes sets a limit of 1 vCPU on containers and 512Mi for memory.

pod-definition.yaml
```
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
  labels:
    app: nginx
spec:
  containers:
    - name: webapp-container
      image: webapp
      resources:
        requests:
          memory: "1Gi"
          cpu: 1
        limits:
          memory: "2Gi"
          cpu: 2
```

If you are exceeding the cpu limit, kubernetes is going to try and throttle the CPU for the container.
However, if the memory usage goes beyond the specified limit kubernetes will terminate the container.
If a pod spec specifies a limit but not a request then the request is set same as the limit values.
However, if a pod spec specifies a request but not a limit, then the limit is set to the limit value mentioned in the `LimitRange` (whether that's custom or default).


---------------------------------------------
The default limits are picked up by the pod in a namespace based on the memory and cpu limit range set in the namespace.
Remember the LimitRange resource is namespace aware and not clusterwide
`cat mem-limit-range.yaml`
```
apiVersion: v1
kind: LimitRange
metadata:
  name: mem-limit-range
spec:
  limits:
  - default:
      memory: 512Mi
    defaultRequest:
      memory: 256Mi
    type: Container
```

AND `cat cpu-limit-range.yaml`
```
apiVersion: v1
kind: LimitRange
metadata:
  name: cpu-limit-range
spec:
  limits:
  - default:
      cpu: 1
    defaultRequest:
      cpu: 0.5
    type: Container
```
