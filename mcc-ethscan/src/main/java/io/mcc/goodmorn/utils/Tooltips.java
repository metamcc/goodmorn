package io.mcc.goodmorn.utils;

import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class Tooltips {

    public static String getOnlyDigits(String s) {
        Pattern pattern = Pattern.compile("[^0-9]");
        Matcher matcher = pattern.matcher(s);
        String number = matcher.replaceAll("");
        return number;
    }
    public static String getOnlyStrings(String s) {
        Pattern pattern = Pattern.compile("[^a-z A-Z]");
        Matcher matcher = pattern.matcher(s);
        String number = matcher.replaceAll("");
        return number;
    }

    public static void main(String[] args) {

        System.out.println(getOnlyDigits("010-6899+;122&5"));
    }

}
