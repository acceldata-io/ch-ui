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
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;

@Controller
public class IndexController {

    @Value("${clickhouse.url}")
    private String clickhouseUrl;

    @Value("${clickhouse.user:}")
    private String clickhouseUser;

    @Value("${clickhouse.pass:}")
    private String clickhousePass;

    @Value("${clickhouse.useAdvanced:false}")
    private String clickhouseUseAdvanced;

    @Value("${clickhouse.customPath:}")
    private String clickhouseCustomPath;

    @Value("${clickhouse.timeout:3000}")
    private String clickhouseTimeout;

    @GetMapping("/")
    public String index(Model model) {
        model.addAttribute("clickhouseUrl", clickhouseUrl);
        model.addAttribute("clickhouseUser", clickhouseUser);
        model.addAttribute("clickhousePass", clickhousePass);
        model.addAttribute("clickhouseUseAdvanced", clickhouseUseAdvanced);
        model.addAttribute("clickhouseCustomPath", clickhouseCustomPath);
        model.addAttribute("clickhouseTimeout", clickhouseTimeout);
        return "index"; // this must be a Thymeleaf template
    }
}