package io.mcc.goodmorn.utils;

import java.util.Hashtable;
import java.util.Locale;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

public class PhoneUtils {
	
	static private Hashtable<String,String> countryList = new Hashtable<String,String>();
	static{
//#define _RELEASE
//#ifdef _RELEASE
		countryList.put("KR","82");
		countryList.put("AD","376");
		countryList.put("AE","971");
		countryList.put("AL","355");
		countryList.put("AM","374");
		countryList.put("AN","599");
		countryList.put("AO","244");
		countryList.put("AR","54");
		countryList.put("AT","43");
		countryList.put("AU","61");
		countryList.put("AW","297");
		countryList.put("AZ","994");
		countryList.put("BA","387");
		countryList.put("BB","1");
		countryList.put("BD","880");
		countryList.put("BE","32");
		countryList.put("BF","226");
		countryList.put("BG","359");
		countryList.put("BH","973");
		countryList.put("BI","257");
		countryList.put("BJ","229");
		countryList.put("BN","673");
		countryList.put("BO","591");
		countryList.put("BR","55");
		countryList.put("BS","1");
		countryList.put("BT","975");
		countryList.put("BW","267");
		countryList.put("BY","375");
		countryList.put("BZ","501");
		countryList.put("CA","1");
		countryList.put("CD","243");
		countryList.put("CF","236");
		countryList.put("CG","242");
		countryList.put("CK","682");
		countryList.put("CL","56");
		countryList.put("CM","237");
		countryList.put("CN","86");
		countryList.put("CO","57");
		countryList.put("CR","506");
		countryList.put("CU","53");
		countryList.put("CV","238");
		countryList.put("CY","357");
		countryList.put("CZ","420");
		countryList.put("DJ","253");
		countryList.put("DO","1");
		countryList.put("DZ","213");
		countryList.put("EC","593");
		countryList.put("EE","372");
		countryList.put("EG","20");
		countryList.put("ER","291");
		countryList.put("ET","251");
		countryList.put("FI","358");
		countryList.put("FJ","679");
		countryList.put("FK","500");
		countryList.put("FM","691");
		countryList.put("FO","298");
		countryList.put("GA","241");
		countryList.put("GD","1");
		countryList.put("GE","995");
		countryList.put("GF","594");
		countryList.put("GH","233");
		countryList.put("GI","350");
		countryList.put("GL","299");
		countryList.put("GM","220");
		countryList.put("GN","224");
		countryList.put("GP","590");
		countryList.put("GQ","240");
		countryList.put("GT","502");
		countryList.put("GU","1");
		countryList.put("GW","245");
		countryList.put("GY","592");
		countryList.put("HK","852");
		countryList.put("HN","504");
		countryList.put("HR","385");
		countryList.put("HT","509");
		countryList.put("HU","36");
		countryList.put("ID","62");
		countryList.put("IL","972");
		countryList.put("IN","91");
		countryList.put("IO","246");
		countryList.put("IQ","964");
		countryList.put("IS","354");
		countryList.put("JM","1");
		countryList.put("JO","962");
		countryList.put("JP","81");
		countryList.put("KE","254");
		countryList.put("KG","996");
		countryList.put("KH","855");
		countryList.put("KI","686");
		countryList.put("KM","269");
		countryList.put("KW","965");
		countryList.put("LA","856");
		countryList.put("LB","961");
		countryList.put("LI","423");
		countryList.put("LK","94");
		countryList.put("LR","231");
		countryList.put("LS","266");
		countryList.put("LT","370");
		countryList.put("LU","352");
		countryList.put("LV","371");
		countryList.put("LY","218");
		countryList.put("MC","377");
		countryList.put("MD","373");
		countryList.put("ME","382");
		countryList.put("MG","261");
		countryList.put("MH","692");
		countryList.put("MK","389");
		countryList.put("ML","223");
		countryList.put("MM","95");
		countryList.put("MN","976");
		countryList.put("MO","853");
		countryList.put("MP","1");
		countryList.put("MQ","596");
		countryList.put("MR","222");
		countryList.put("MT","356");
		countryList.put("MU","230");
		countryList.put("MV","960");
		countryList.put("MW","265");
		countryList.put("MX","52");
		countryList.put("MY","60");
		countryList.put("MZ","258");
		countryList.put("NA","264");
		countryList.put("NC","687");
		countryList.put("NE","227");
		countryList.put("NG","234");
		countryList.put("NI","505");
		countryList.put("NR","674");
		countryList.put("NU","683");
		countryList.put("NZ","64");
		countryList.put("PA","507");
		countryList.put("PE","51");
		countryList.put("PF","689");
		countryList.put("PG","675");
		countryList.put("PH","63");
		countryList.put("PK","92");
		countryList.put("PM","508");
		countryList.put("PR","1");
		countryList.put("PW","680");
		countryList.put("PY","595");
		countryList.put("RE","262");
		countryList.put("RO","40");
		countryList.put("RS","381");
		countryList.put("SA","966");
		countryList.put("SB","677");
		countryList.put("SC","248");
		countryList.put("SD","249");
		countryList.put("SE","46");
		countryList.put("SG","65");
		countryList.put("SH","290");
		countryList.put("SI","386");
		countryList.put("SK","421");
		countryList.put("SL","232");
		countryList.put("SM","378");
		countryList.put("SN","221");
		countryList.put("SO","252");
		countryList.put("SR","597");
		countryList.put("ST","239");
		countryList.put("SV","503");
		countryList.put("SY","963");
		countryList.put("SZ","268");
		countryList.put("TD","235");
		countryList.put("TG","228");
		countryList.put("TH","66");
		countryList.put("TJ","992");
		countryList.put("TK","690");
		countryList.put("TL","670");
		countryList.put("TM","993");
		countryList.put("TW","886");
		countryList.put("TZ","255");
		countryList.put("UA","380");
		countryList.put("US","1");
		countryList.put("UY","598");
		countryList.put("UZ","998");
		countryList.put("VC","1");
		countryList.put("VE","58");
		countryList.put("VN","84");
		countryList.put("VU","678");
		countryList.put("WF","681");
		countryList.put("WS","685");
		countryList.put("YE","967");
		countryList.put("ZM","260");
		countryList.put("ZW","263");
	//#else	
//@		countryList.put("AC","247");
//@		countryList.put("AD","376");
//@		countryList.put("AE","971");
//@		countryList.put("AF","93");
//@		countryList.put("AG","1");
//@		countryList.put("AI","1");
//@		countryList.put("AL","355");
//@		countryList.put("AM","374");
//@		countryList.put("AN","599");
//@		countryList.put("AO","244");
//@		countryList.put("AQ","672");
//@		countryList.put("AR","54");
//@		countryList.put("AS","1");
//@		countryList.put("AT","43");
//@		countryList.put("AU","61");
//@		countryList.put("AW","297");
//@		countryList.put("AZ","994");
//@		countryList.put("BA","387");
//@		countryList.put("BB","1");
//@		countryList.put("BD","880");
//@		countryList.put("BE","32");
//@		countryList.put("BF","226");
//@		countryList.put("BG","359");
//@		countryList.put("BH","973");
//@		countryList.put("BI","257");
//@		countryList.put("BJ","229");
//@		countryList.put("BL","590");
//@		countryList.put("BM","1");
//@		countryList.put("BN","673");
//@		countryList.put("BO","591");
//@		countryList.put("BR","55");
//@		countryList.put("BS","1");
//@		countryList.put("BT","975");
//@		countryList.put("BW","267");
//@		countryList.put("BY","375");
//@		countryList.put("BZ","501");
//@		countryList.put("CA","1");
//@		countryList.put("CC","61");
//@		countryList.put("CD","243");
//@		countryList.put("CF","236");
//@		countryList.put("CG","242");
//@		countryList.put("CH","41");
//@		countryList.put("CI","225");
//@		countryList.put("CK","682");
//@		countryList.put("CL","56");
//@		countryList.put("CM","237");
//@		countryList.put("CN","86");
//@		countryList.put("CO","57");
//@		countryList.put("CR","506");
//@		countryList.put("CU","53");
//@		countryList.put("CV","238");
//@		countryList.put("CW","599");
//@		countryList.put("CX","61");
//@		countryList.put("CY","357");
//@		countryList.put("CZ","420");
//@		countryList.put("DE","49");
//@		countryList.put("DJ","253");
//@		countryList.put("DK","45");
//@		countryList.put("DM","1");
//@		countryList.put("DO","1");
//@		countryList.put("DZ","213");
//@		countryList.put("EC","593");
//@		countryList.put("EE","372");
//@		countryList.put("EG","20");
//@		countryList.put("EH","212");
//@		countryList.put("ER","291");
//@		countryList.put("ES","34");
//@		countryList.put("ET","251");
//@		countryList.put("FI","358");
//@		countryList.put("FJ","679");
//@		countryList.put("FK","500");
//@		countryList.put("FM","691");
//@		countryList.put("FO","298");
//@		countryList.put("FR","33");
//@		countryList.put("GA","241");
//@		countryList.put("GB","44");
//@		countryList.put("GD","1");
//@		countryList.put("GE","995");
//@		countryList.put("GF","594");
//@		countryList.put("GG","44");
//@		countryList.put("GH","233");
//@		countryList.put("GI","350");
//@		countryList.put("GL","299");
//@		countryList.put("GM","220");
//@		countryList.put("GN","224");
//@		countryList.put("GP","590");
//@		countryList.put("GQ","240");
//@		countryList.put("GR","30");
//@		countryList.put("GT","502");
//@		countryList.put("GU","1");
//@		countryList.put("GW","245");
//@		countryList.put("GY","592");
//@		countryList.put("HK","852");
//@		countryList.put("HN","504");
//@		countryList.put("HR","385");
//@		countryList.put("HT","509");
//@		countryList.put("HU","36");
//@		countryList.put("ID","62");
//@		countryList.put("IE","353");
//@		countryList.put("IL","972");
//@		countryList.put("IM","44");
//@		countryList.put("IN","91");
//@		countryList.put("IO","246");
//@		countryList.put("IQ","964");
//@		countryList.put("IR","98");
//@		countryList.put("IS","354");
//@		countryList.put("IT","39");
//@		countryList.put("JE","44");
//@		countryList.put("JM","1");
//@		countryList.put("JO","962");
//@		countryList.put("JP","81");
//@		countryList.put("KE","254");
//@		countryList.put("KG","996");
//@		countryList.put("KH","855");
//@		countryList.put("KI","686");
//@		countryList.put("KM","269");
//@		countryList.put("KN","1");
//@		countryList.put("KP","850");
//@		countryList.put("KR","82");
//@		countryList.put("KW","965");
//@		countryList.put("KY","1");
//@		countryList.put("KZ","7");
//@		countryList.put("LA","856");
//@		countryList.put("LB","961");
//@		countryList.put("LC","1");
//@		countryList.put("LI","423");
//@		countryList.put("LK","94");
//@		countryList.put("LR","231");
//@		countryList.put("LS","266");
//@		countryList.put("LT","370");
//@		countryList.put("LU","352");
//@		countryList.put("LV","371");
//@		countryList.put("LY","218");
//@		countryList.put("MA","212");
//@		countryList.put("MC","377");
//@		countryList.put("MD","373");
//@		countryList.put("ME","382");
//@		countryList.put("MF","1");
//@		countryList.put("MG","261");
//@		countryList.put("MH","692");
//@		countryList.put("MK","389");
//@		countryList.put("ML","223");
//@		countryList.put("MM","95");
//@		countryList.put("MN","976");
//@		countryList.put("MO","853");
//@		countryList.put("MP","1");
//@		countryList.put("MQ","596");
//@		countryList.put("MR","222");
//@		countryList.put("MS","1");
//@		countryList.put("MT","356");
//@		countryList.put("MU","230");
//@		countryList.put("MV","960");
//@		countryList.put("MW","265");
//@		countryList.put("MX","52");
//@		countryList.put("MY","60");
//@		countryList.put("MZ","258");
//@		countryList.put("NA","264");
//@		countryList.put("NC","687");
//@		countryList.put("NE","227");
//@		countryList.put("NF","672");
//@		countryList.put("NG","234");
//@		countryList.put("NI","505");
//@		countryList.put("NL","31");
//@		countryList.put("NO","47");
//@		countryList.put("NP","977");
//@		countryList.put("NR","674");
//@		countryList.put("NU","683");
//@		countryList.put("NZ","64");
//@		countryList.put("OM","968");
//@		countryList.put("PA","507");
//@		countryList.put("PE","51");
//@		countryList.put("PF","689");
//@		countryList.put("PG","675");
//@		countryList.put("PH","63");
//@		countryList.put("PK","92");
//@		countryList.put("PL","48");
//@		countryList.put("PM","508");
//@		countryList.put("PR","1");
//@		countryList.put("PT","351");
//@		countryList.put("PW","680");
//@		countryList.put("PY","595");
//@		countryList.put("QA","974");
//@		countryList.put("RE","262");
//@		countryList.put("RO","40");
//@		countryList.put("RS","381");
//@		countryList.put("RU","7");
//@		countryList.put("RW","250");
//@		countryList.put("SA","966");
//@		countryList.put("SB","677");
//@		countryList.put("SC","248");
//@		countryList.put("SD","249");
//@		countryList.put("SE","46");
//@		countryList.put("SG","65");
//@		countryList.put("SH","290");
//@		countryList.put("SI","386");
//@		countryList.put("SJ","47");
//@		countryList.put("SK","421");
//@		countryList.put("SL","232");
//@		countryList.put("SM","378");
//@		countryList.put("SN","221");
//@		countryList.put("SO","252");
//@		countryList.put("SR","597");
//@		countryList.put("SS","211");
//@		countryList.put("ST","239");
//@		countryList.put("SV","503");
//@		countryList.put("SX","1");
//@		countryList.put("SY","963");
//@		countryList.put("SZ","268");
//@		countryList.put("TC","1");
//@		countryList.put("TD","235");
//@		countryList.put("TG","228");
//@		countryList.put("TH","66");
//@		countryList.put("TJ","992");
//@		countryList.put("TK","690");
//@		countryList.put("TL","670");
//@		countryList.put("TM","993");
//@		countryList.put("TN","216");
//@		countryList.put("TO","676");
//@		countryList.put("TR","90");
//@		countryList.put("TT","1");
//@		countryList.put("TV","688");
//@		countryList.put("TW","886");
//@		countryList.put("TZ","255");
//@		countryList.put("UA","380");
//@		countryList.put("UG","256");
//@		countryList.put("US","1");
//@		countryList.put("UY","598");
//@		countryList.put("UZ","998");
//@		countryList.put("VA","379");
//@		countryList.put("VC","1");
//@		countryList.put("VE","58");
//@		countryList.put("VG","1");
//@		countryList.put("VI","1");
//@		countryList.put("VN","84");
//@		countryList.put("VU","678");
//@		countryList.put("WF","681");
//@		countryList.put("WS","685");
//@		countryList.put("YE","967");
//@		countryList.put("YT","262");
//@		countryList.put("ZA","27");
//@		countryList.put("ZM","260");
//@		countryList.put("ZW","263");
	//#endif	
	}
	static public String getCallingCode(String country_code){
		String ccode_2digit = country_code;
		if (country_code.length()==3)  {
		    Locale locale = new Locale("en",country_code);
		    System.out.println("Country=" + locale.getISO3Country());
		    
			ccode_2digit = country_code.substring(0, 2);
		}
		if(countryList.containsKey(ccode_2digit)){
			return countryList.get(ccode_2digit);
		}else{
			return "";
		}
	}
	
	public static String convertInphoneToOutphone(String shphone, String countryCode){
		String ret = shphone;
		String callingCode = getCallingCode(countryCode);
		//if("KR".equals(countryCode)){
			if(shphone != null && shphone.length() > 2){
				if("0".equals(shphone.substring(0,1)) && !"0".equals(shphone.substring(1,2))){
					if(!"".equals(callingCode)){
						ret = callingCode + shphone.substring(1);
					}
				}else if("82".equals(callingCode) && shphone.startsWith("999")){	//테스트용번호
					
				}else if(!"0".equals(shphone.substring(0,1)) && !shphone.startsWith(callingCode)){	//sinyu.2015.03.09 시작이 0이아니고 국가코드로 시작하지 않으면 전화번호앞에 국가코드를 삽입(미국,중국등은 0을 앞에 사용하지 않느다)
					ret = callingCode + shphone.substring(0);
				}
			}
		//}
		return ret;
	}
	public static String convertOutphoneToInphone(String shphone, String countryCode){
		String ret = shphone;
		String callingCode = getCallingCode(countryCode);
		if(!"".equals(callingCode)){
			String callphone = shphone;
			if (shphone.startsWith("+")) callphone = shphone.substring(1);
			
			if(callphone != null && callphone.length() > callingCode.length()){
				if(callphone.startsWith(callingCode)){	//국제코드로 시작하면 국제코드를 제거
					ret = callphone.substring(callingCode.length());
					if("81".equals(callingCode) || "82".equals(callingCode)){	//한국/일본이면 0을추가
						if (ret.charAt(0) != '0')
							ret = "0" + ret;
						else return ret;
					}
				}else if("82".equals(callingCode) && shphone.startsWith("999")){	//테스트용번호
					
				}else if(("81".equals(callingCode) || "82".equals(callingCode)) && !shphone.startsWith("0")){	//국제코드로 시작않하고 한국/일본이면 0을 추가
					ret = "0" + ret;
				}
			}
		} 
		return ret;
	}

	private static final String PHONE_NUM_PATTERN = "(01[016789])(\\d{3,4})(\\d{4})";
	private static boolean isValid(final String regex, final String target) {
		Matcher matcher = Pattern.compile(regex).matcher(target);
		return matcher.matches();
	}
	public static boolean isPhoneNum(final String str) {
		return isValid(PHONE_NUM_PATTERN, str);
	}

}
