package com.k8s_demo;

import org.apache.commons.cli.*;

import io.grpc.examples.helloworld.*;
import io.grpc.Channel;
import io.grpc.Grpc;
import io.grpc.InsecureChannelCredentials;
import io.grpc.ManagedChannel;
import io.grpc.StatusRuntimeException;
import io.grpc.TlsChannelCredentials;

import java.util.concurrent.TimeUnit;
import java.util.logging.Level;
import java.util.logging.Logger;
import java.io.File;

/**
 * A simple client that requests a greeting from the {@link HelloWorldServer}.
 */
public class HelloWorldClient {
    private static final Logger logger = Logger.getLogger(HelloWorldClient.class.getName());

    private final GreeterGrpc.GreeterBlockingStub blockingStub;

    /** Construct client for accessing HelloWorld server using the existing channel. */
    public HelloWorldClient(Channel channel) {
        // 'channel' here is a Channel, not a ManagedChannel, so it is not this code's responsibility to
        // shut it down.

        // Passing Channels to code makes code easier to test and makes it easier to reuse Channels.
        blockingStub = GreeterGrpc.newBlockingStub(channel);
    }

    /** Say hello to server. */
    public void greet(String name) {
        logger.info("Will try to greet " + name + " ...");
        HelloRequest request = HelloRequest.newBuilder().setName(name).build();
        HelloReply response;
        try {
            response = blockingStub.sayHello(request);
        } catch (StatusRuntimeException e) {
            logger.log(Level.WARNING, "RPC failed: {0}", e.getStatus());
            return;
        }
        logger.info("Greeting: " + response.getMessage());
    }

    /**
     * Greet server. If provided, the first element of {@code args} is the name to use in the
     * greeting. The second argument is the target server.
     */
    public static void main(String[] args) throws Exception {
        Options options = new Options();
        options.addOption(Option.builder("target")
                .required() // Make it required
                .hasArg()
                .desc("url")
                .build());
        options.addOption(Option.builder("caKey")
                .required() // Make it required
                .hasArg()
                .desc("Certificating Authority")
                .build());
        options.addOption(Option.builder("Cert")
                .required() // Make it required
                .hasArg()
                .desc("Certificating Authority")
                .build());

        options.addOption(Option.builder("Key")
                .required() // Make it required
                .hasArg()
                .desc("Certificating Authority")
                .build());


        CommandLineParser parser = new DefaultParser();
        CommandLine cmd = parser.parse(options, args);


        // Create a communication channel to the server, known as a Channel. Channels are thread-safe
        // and reusable. It is common to create channels at the beginning of your application and reuse
        // them until the application shuts down.
        //
        // For the example we use plaintext insecure credentials to avoid needing TLS certificates. To
        // use TLS, use TlsChannelCredentials instead.
        TlsChannelCredentials.Builder tlsBuilder = TlsChannelCredentials.newBuilder();
        tlsBuilder.keyManager(new File(cmd.getOptionValue("Cert")), new File(cmd.getOptionValue("Key")));
        tlsBuilder.trustManager(new File((cmd.getOptionValue("caKey"))));

        while (true) {
            ManagedChannel channel = Grpc.newChannelBuilder(cmd.getOptionValue("target"), tlsBuilder.build())
                    .build();
            for (int i = 0; i< 100; i++) {
                try {
                    HelloWorldClient client = new HelloWorldClient(channel);
                    client.greet("user");
                } finally {
                    // ManagedChannels use resources like threads and TCP connections. To prevent leaking these
                    // resources the channel should be shut down when it will no longer be used. If it may be used
                    // again leave it running.
                    channel.shutdownNow().awaitTermination(5, TimeUnit.SECONDS);
                }
                TimeUnit.SECONDS.sleep(1);
            }
        }
    }
}
