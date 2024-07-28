# go-nats-sample

The sample about NATS, which use these libraries:
- [nats-io/nats.go](https://github.com/nats-io/nats.go)
- [core-go/nats](https://github.com/core-go/nats) to wrap [nats-io/nats.go](https://github.com/nats-io/nats.go)
    - Simplify the way to initialize the consumer, publisher by configurations
        - Props: when you want to change the parameter of consumer or publisher, you can change the config file, and restart Kubernetes POD, do not need to change source code and re-compile.
- [core-go/mq](https://github.com/core-go/mq) to implement this flow, which can be considered a low code tool for message queue consumer:

  ![A common flow to consume a message from a message queue](https://cdn-images-1.medium.com/max/800/1*Y4QUN6QnfmJgaKigcNHbQA.png)

### Similar libraries for nodejs
We also provide these libraries to support nodejs:
- [nats-plus](https://www.npmjs.com/package/nats-plus), to wrap and simplify [nats](https://www.npmjs.com/package/nats). The sample is at [nats-sample](https://github.com/typescript-tutorial/nats-sample)
