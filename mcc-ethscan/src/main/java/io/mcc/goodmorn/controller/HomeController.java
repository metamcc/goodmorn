package io.mcc.goodmorn.controller;

import io.mcc.common.entity.EthMccTxLog;
import io.mcc.common.entity.EthTransRes;
import io.mcc.goodmorn.service.MCCTxService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@Slf4j
public class HomeController {

    @Autowired
    private MCCTxService mccTxService;

    @RequestMapping("/")
    public @ResponseBody  String welcome(Model model) {

        return "Hello! mcc ethscan";
    }
    @RequestMapping(value = "/api", method = RequestMethod.GET)
    public @ResponseBody EthTransRes mccScanApi(
            @RequestParam("startblock") Long startblock,
            @RequestParam("endblock") Long endblock,
            @RequestParam("sort") String orderby) {

        EthTransRes txRes = mccTxService.findMccTxLog(startblock, endblock, orderby);
        return txRes;
    }
}
