package io.mcc.auth.common.util;

import com.google.common.base.Stopwatch;
import org.apache.commons.lang3.StringUtils;
import org.slf4j.MDC;
import org.springframework.web.servlet.handler.HandlerInterceptorAdapter;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.net.Inet4Address;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.UUID;

public class GenerateLogTransactionKey extends HandlerInterceptorAdapter {
    private DateTimeFormatter formatter = DateTimeFormatter.ofPattern("yyyyMMddHHmmss");

    @Override
    public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler) throws Exception {
        final String ipAddress[] = Inet4Address.getLocalHost().getHostAddress().split("[.]");
        final String postIp = StringUtils.leftPad(ipAddress[ipAddress.length-1], 3, "0");
        final String tranKey = postIp + LocalDateTime.now().format(formatter)
                + UUID.randomUUID().toString().replaceAll("-", "");

        MDC.put("tranKey", tranKey);
        final Stopwatch watch = Stopwatch.createStarted();
        request.setAttribute("stopWatch", watch);
        return super.preHandle(request, response, handler);
    }
}
