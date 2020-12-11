Rate-limiter used for restrict high-concurrency requests. Common algorithms inclued counter, token bucket and leaky bucket. In this repo, I use thoes three algorithms to implement a rate-limiter which also supports instant large amounts of requests.

#### Counter

Easiest way to implement a rate-limiter. We just need to judge whether the number of requests exceeds the limit. We also need a timer to reset requests in windows. The difference of a fixed window and a sliding window is the granularity of the step the window moves.

1. Fixed window

   Suspect QPS is 1k/s, we need to judge whether the counter reaches 1k, if so reject request, otherwise execute the request.

2. Sliding window

   The sliding window can be divided into a more fine-grained based on the basis of the fixed window, such as dividing the 1s request into 1000 buckets, and each bucket represents a time range of 1ms. Then it becomes a fixed window. The code I implement is a little different from the fixed window, because I was just trying to practicing writing golang. haha=.= 

