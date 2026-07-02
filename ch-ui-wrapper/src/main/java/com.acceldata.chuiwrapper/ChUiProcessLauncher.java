/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package com.acceldata.chuiwrapper;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.CommandLineRunner;
import org.springframework.stereotype.Component;

import java.io.File;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.TimeUnit;

/**
 * CH-UI ships as a single Go binary that already serves the frontend, the
 * API, and the websocket tunnel on its own. This wrapper does not
 * re-implement any of that; it only launches and supervises the Go
 * "ch-ui server" process so it can be managed as a JVM service (PID file,
 * start/stop lifecycle) by ODP tooling.
 */
@Component
public class ChUiProcessLauncher implements CommandLineRunner {

    @Value("${chui.binaryPath:../ch-ui}")
    private String binaryPath;

    @Value("${chui.port:3488}")
    private int port;

    @Value("${clickhouse.url:http://localhost:8123}")
    private String clickhouseUrl;

    @Value("${clickhouse.connectionName:Local ClickHouse}")
    private String connectionName;

    @Value("${chui.pidFile:ch-ui-server.pid}")
    private String childPidFile;

    @Value("${chui.stopTimeoutSeconds:15}")
    private long stopTimeoutSeconds;

    private Process childProcess;

    @Override
    public void run(String... args) throws Exception {
        File binary = new File(binaryPath);
        if (!binary.isFile() || !binary.canExecute()) {
            throw new IllegalStateException(
                    "ch-ui binary not found or not executable at " + binary.getAbsolutePath()
                            + ". Set chui.binaryPath to the compiled Go ch-ui executable.");
        }

        List<String> command = new ArrayList<>();
        command.add(binary.getAbsolutePath());
        command.add("server");
        command.add("--port");
        command.add(String.valueOf(port));
        if (clickhouseUrl != null && !clickhouseUrl.isBlank()) {
            command.add("--clickhouse-url");
            command.add(clickhouseUrl);
        }
        if (connectionName != null && !connectionName.isBlank()) {
            command.add("--connection-name");
            command.add(connectionName);
        }
        command.add("--pid-file");
        command.add(childPidFile);

        ProcessBuilder builder = new ProcessBuilder(command)
                .redirectOutput(ProcessBuilder.Redirect.INHERIT)
                .redirectError(ProcessBuilder.Redirect.INHERIT);

        childProcess = builder.start();
        Runtime.getRuntime().addShutdownHook(new Thread(this::stopChild));

        int exitCode = childProcess.waitFor();
        if (exitCode != 0) {
            System.exit(exitCode);
        }
    }

    private void stopChild() {
        if (childProcess == null || !childProcess.isAlive()) {
            return;
        }
        childProcess.destroy();
        try {
            if (!childProcess.waitFor(stopTimeoutSeconds, TimeUnit.SECONDS)) {
                childProcess.destroyForcibly();
            }
        } catch (InterruptedException e) {
            childProcess.destroyForcibly();
            Thread.currentThread().interrupt();
        }
    }
}
