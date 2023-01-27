package io.mcc.auth.common.util;

import java.time.Instant;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.LocalTime;
import java.time.ZoneId;
import java.time.format.DateTimeFormatter;
import java.util.ArrayList;
import java.util.Calendar;
import java.util.Collections;
import java.util.Enumeration;
import java.util.HashSet;
import java.util.Hashtable;
import java.util.List;
import java.util.Map;
import java.util.Random;
import java.util.Set;

import javax.servlet.ServletRequest;

import lombok.extern.slf4j.Slf4j;

@Slf4j
public class Util
{
	public static final String SUCCESS_KEY  =  "success";
	public static final String MSG_CODE_KEY =  "msg_code";
	public static final String DATA_KEY     =  "data";
	public static final String MSG_KEY      =  "msg";
	public static final String PACKET_DATA  =  "packet_data";
	
	public static Map<String, Object> getReqToMap(ServletRequest req){
		Map<String, Object>     model = new Hashtable<String, Object>();
		Enumeration<String> enu = req.getParameterNames();
		
		String key = null;
		while(enu.hasMoreElements()){
			key = enu.nextElement();
			if(key != null){
				model.put(key, req.getParameter(key));
			}
		}
		
		return model;
	}
	
	public static int toInt(String count)
	{
		try
		{
			return Integer.parseInt(count.trim());
		}
		catch(Exception e)
		{
			return 0;
		}
	}

	public static int objToInt(Object obj){
		if(obj == null){
			return 0;
		} else {
			try
			{
				return Integer.parseInt(String.valueOf(obj));
			}
			catch(Exception e)
			{
				return 0;
			}
		}
	}
	
	public static long objToLong(Object obj){
		if(obj == null){
			return 0;
		} else {
			try
			{
				return Long.parseLong(String.valueOf(obj));
			}
			catch(Exception e)
			{
				return 0;
			}
		}		
	}
	
	public static float objToFloat(Object obj){
		if(obj == null){
			return 0;
		} else {
			try
			{
				return Float.parseFloat(String.valueOf(obj));
			}
			catch(Exception e)
			{
				return 0;
			}
		}
	}
	
	public static String nullTobk(String str){
		
		if(str == null){
			return "";
		}else if("null".equals(str.toLowerCase())){
			return "";
		} else {
			return str;
		}
	}
	
	public static String objToStr(Object obj){
	
		if(obj == null){
			return "";
		}else{
			String s = String.valueOf(obj);
			if("null".equals(s.toLowerCase())){
					return "";
			} else {
				return s;
			}
		}
	}
	
	public static boolean objToBoolean(Object obj){
		if(obj == null){
			return false;
		} else {
			try
			{
				return Boolean.parseBoolean(String.valueOf(obj));
			}
			catch(Exception e)
			{
				return false;
			}
		}
	}
	
	public static String toDayText(long toDay)
	{
		StringBuffer buffer = new StringBuffer();
		
		Calendar calendar = Calendar .getInstance();
		
		calendar.setTimeInMillis(toDay);
		
		buffer.append((calendar.get(Calendar.YEAR))+"/");
		buffer.append((calendar.get(Calendar.MONTH)+1+"/"));
		buffer.append((calendar.get(Calendar.DATE))+" ");
		buffer.append(calendar.get(Calendar.HOUR)+":");
		buffer.append(calendar.get(Calendar.MINUTE)+":");
		buffer.append(calendar.get(Calendar.SECOND));

		return buffer.toString();
	}
	
	public static String toDay(String zi)
	{
		LocalDateTime userDateTime = LocalDateTime.now(ZoneId.of(zi));
		
		String format_userDate = userDateTime.format(DateTimeFormatter.ofPattern("yyyy년MM월dd일"));
		
		return format_userDate;
	}

	public static String toDay(String zi, int addDay)
	{
		LocalDate userDateTime = LocalDate.now(ZoneId.of(zi));
		LocalDate valuTime = userDateTime.plusDays(addDay);
		String format_userDate = valuTime.format(DateTimeFormatter.ofPattern("yyyy년MM월dd일"));
		
		return format_userDate;
	}
	
	public static boolean isAfterDay(String endDate, String  zi) {
		LocalDate end   = LocalDate.parse(endDate);
		LocalDate start = LocalDate.now(ZoneId.of(zi));
		// start이 end 보다 이후 날짜인지 비교
		return start.isAfter(end);
		
	}
	
	public static boolean setProbability(int random_count , int count){
		Set<Integer> list = new HashSet<Integer>();
		Random oRandom = new Random();
		boolean k = false;
		for(int i=0; i<count; i++){
			list.add((i+1));
		}
		
		int i = oRandom.nextInt(list.size());
		Object[] arr = list.toArray();
		
		i = Util.objToInt(arr[i]);
		list.clear();
		
		if( i < random_count){
			k = true;
		}
		
		return k;
	}
	
	public static Object setItemProbability(String pKey, Object[] key, Map<String, Map<String,Object>> boost){
		List<Object> list = new ArrayList<Object>();
		Object o = null;
		Map<String, Object> map = null;
		Random oRandom = new Random();
		int k = 0;
		int count  = key.length;
		
		for(int i=0; i<count; i++){
			map = boost.get(key[i]);
			k = Util.objToInt(map.get(pKey));
			k = k*10;
			for(int j = 0; j<k; j++){
				list.add(key[i]);
			}
		}
		
		Collections.shuffle(list);
		
		int i = oRandom.nextInt((list.size()-1));
		o = list.get(i);
		list.clear();
		return o;
	}
	
	@SuppressWarnings("unchecked")
	public static void itemDataGroupSet(List<Object> list, String group_id, Map<String, List<Map<String, Object>>> store){
		Map<String, Object> map            = null;
		List<Map<String, Object>>  data    = null;
		int size = list.size();

		for(int i = 0; i < size; i++){
			map = (Map<String, Object>)list.get(i);
			String item_group_id = Util.objToStr(map.get(group_id));
			data = store.get(item_group_id);
			if(data == null || data.size() <  1){
				data = new ArrayList<Map<String, Object>>(); 
				data.add(map);
			} else {
				data.add(map);
			}
			
			store.put(item_group_id, data);
		}
	}
	
	
	public static Object setListProbability(List<Map<String, Object>> list, String key1, String key2){
		return setListProbability(list, key1, key2, 0);
	}
	
	public static Object setListProbability(List<Map<String, Object>> list, String key1, String key2, int d){
		List<Object> p = new ArrayList<Object>();
		int count = list.size();
		Map<String, Object> map = null;
		Object o = null;
		Random oRandom = new Random();
		int k = 0;
		for(int i=0; i<count; i++){
			map = list.get(i);
			k = Util.objToInt(map.get(key1));
			k   =  k  <1  ? d : k;
			k = k*10;
			for(int j = 0; j<k; j++){
				p.add(map.get(key2));
			}
		}
		
		Collections.shuffle(p);
		try{
			int i = oRandom.nextInt((p.size()-1));
			
			if(i < 0 || i > (p.size()-1)){
				i =0;
			}
			
			o = p.get(i);
		}catch(Exception e){
			o = p.get(0);
		}
		
		p.clear();
		return o;
		
	}
	
	public static boolean setProbabilityBoolean(int count, int clearCount){
		List<Boolean> list = new ArrayList<Boolean>();
		Random oRandom = new Random();
		for(int i=0; i<count; i++){
			if(i < clearCount){
				list.add(true);
			} else {
				list.add(false);
			}
		}
		
		Collections.shuffle(list);
		int i = oRandom.nextInt(list.size());
		Boolean k =list.get(i);
		list.clear();
		return k;
	}
	
	public static int setProbabilityInt(int count, int maxNum){
		List<Integer> list = new ArrayList<Integer>();
		Random oRandom = new Random();
		int num = 0;
		
		for(int i=0; i<count; i++){
			if(num > maxNum){
				num = 0;
			}
			
			list.add(num);
			num = num +1;
		}
		
		Collections.shuffle(list);
		System.out.println("list   :   "+list);
		int i = oRandom.nextInt(list.size());
		Integer k =list.get(i);
		list.clear();
		return k;
	}
	
	public static String setOsType(Map<String, Object>  json_packet){
		String os_type = Util.objToStr(json_packet.get("os_type"));
		os_type  = os_type.toLowerCase();

		if("android".equals(os_type) || "a".equals(os_type) || "google".equals(os_type)){
			return "A";
		} else if("ios".equals(os_type) || "i".equals(os_type)){
			return "I";
		} else if("window".equals(os_type)){
			return "W";
		} else {
			return "U";
		}
	}
	
	public static String bytesToHex(byte[] data)
	{
		if (data==null) {
			return null;
		}

		int len = data.length;
		String str = "";

		   // 0x는 16진수를 나타낸다.
		   // 0xff는 10진수로 치면 255
		   // 1111 1111 => 128 + 64 + 32 +16 + 8 + 4 + 2 + 1 = 255
		for (int i=0; i<len; i++) 
		{
			if ((data[i]&0xFF)<16) { // 2자리 포맷 맞춰주기
				   str = str + "0" + Integer.toHexString(data[i]&0xFF);
			} else {
				str = str + Integer.toHexString(data[i]&0xFF);
			}
		}
	
		return str;
	}

	public static byte[] hexToBytes(String str) 
	{
		if (str==null) {
			return null;
		} else if (str.length() < 2) {
			return null;
		} 
		else 
		{
			int len = str.length() / 2;
			byte[] buffer = new byte[len];
	   
			for (int i=0; i<len; i++) 
			{
				// 16진수 이므로 숫자 변환시 radix를 명시 한다.
				buffer[i] = (byte) Integer.parseInt(str.substring(i*2,i*2+2),16);
			}
			
			return buffer;
		}
	}

	public static  Map<String, Object> setSuccessCode(Map<String, Object> model){
    	model.put(MSG_CODE_KEY, "000");
    	model.put(SUCCESS_KEY,   true);
        //model.put(DATA_KEY, "");
        model.put(MSG_KEY, "");
        return model;
    }
	
	public static  Map<String, Object> setDefaultfailCode(){
		Map<String, Object>     model = new Hashtable<String, Object>();
		model.put(MSG_CODE_KEY, "err_goods_000");
	    model.put(MSG_KEY,      "");
	    model.put(SUCCESS_KEY,  false);
	    return model;
    }
	
    public static void setBlockNum(int curPage, int block_scale, int totPage, Map<String, String>  model){
    	curPage = curPage < 1 ? 1:curPage;
    	block_scale = block_scale < 5 ? 5:block_scale;
    	// *현재 페이지가 몇번째 페이지 블록에 속하는지 계산
        int curBlock = (int)Math.ceil((curPage-1) / block_scale)+1;
        // *현재 페이지 블록의 시작, 끝 번호 계산
        int blockBegin = (curBlock-1)*block_scale+1;
        // 페이지 블록의 끝번호
        int blockEnd = blockBegin+block_scale-1;
        // *마지막 블록이 범위를 초과하지 않도록 계산
        if(blockEnd > totPage) blockEnd = totPage;
        
        System.out.println("blockEnd  ==>>"+blockEnd);
        System.out.println("blockBegin  ==>>"+blockBegin);
        System.out.println("block_scale  ==>>"+block_scale);
        System.out.println("totPage  ==>>"+totPage);
        
        model.put("blockBegin", objToStr(blockBegin));
        model.put("blockEnd",   objToStr(blockEnd));
    }
	
	
	public static Map<String, Object> getRangeDate(String zoneId) {
		Map<String, Object>     model = new Hashtable<String, Object>();
		log.info("getRangeDate zoneId ==>> "+zoneId);
		ZoneId                  zi = ZoneId.of(zoneId);
		log.info("getRangeDate zi ==>> "+zi);
		LocalDateTime userDateTime = LocalDateTime.ofInstant(Instant.now(), zi);		
		LocalDateTime startOfDate = userDateTime.with(LocalTime.MIN);//2018-10-22T00:00;
		LocalDateTime endOfDate = userDateTime.with(LocalTime.MAX);//2018-10-22T23:59:59.999999999
		log.info("getRangeDate endOfDate ==>> "+endOfDate);
		String format_startOfDate = startOfDate.format(DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss"));
		String format_endOfDate = endOfDate.format(DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss"));
		
		model.put("startOfDate", format_startOfDate);
		model.put("endOfDate", format_endOfDate);
		
		return model;
	}

}