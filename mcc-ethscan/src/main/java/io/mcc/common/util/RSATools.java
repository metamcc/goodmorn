package io.mcc.common.util;

import javax.crypto.Cipher;
import java.io.InputStream;
import java.security.*;
import java.util.Base64;

import static java.nio.charset.StandardCharsets.UTF_8;

public class RSATools {
    public static KeyPair generateKeyPair() throws Exception {
        KeyPairGenerator generator = KeyPairGenerator.getInstance("RSA");
        generator.initialize(2048, new SecureRandom());
        KeyPair pair = generator.generateKeyPair();

        return pair;
    }

    public static KeyPair getKeyPairFromKeyStore(String keystorePath, String keystorePW, String keyAlias) throws Exception {
        //Generated with:
        //keytool -genkey -v -keystore mcc_batch.jks -alias batch_private -keyalg RSA -sigalg MD5withRSA -keysize 2048 -validity 3650

        InputStream ins = RSATools.class.getResourceAsStream(keystorePath);

        //System.out.println("Ins:"+ins);

        KeyStore keyStore = KeyStore.getInstance(KeyStore.getDefaultType());    //"JCEKS");
        keyStore.load(ins, keystorePW.toCharArray());   //Keystore password
        KeyStore.PasswordProtection keyPassword =       //Key password
                new KeyStore.PasswordProtection(keystorePW.toCharArray());

        KeyStore.PrivateKeyEntry privateKeyEntry = (KeyStore.PrivateKeyEntry) keyStore.getEntry(keyAlias, keyPassword);

        java.security.cert.Certificate cert = keyStore.getCertificate(keyAlias);
        PublicKey publicKey = cert.getPublicKey();
        PrivateKey privateKey = privateKeyEntry.getPrivateKey();

        return new KeyPair(publicKey, privateKey);
    }

    public static KeyPair getKeyPairFromKeyStore(InputStream keystoreIn, String keystorePW, String keyAlias) throws Exception {
        //Generated with:
        //keytool -genkey -v -keystore mcc_batch.jks -alias batch_private -keyalg RSA -sigalg MD5withRSA -keysize 2048 -validity 3650


        //System.out.println("Ins:"+ins);

        KeyStore keyStore = KeyStore.getInstance(KeyStore.getDefaultType());    //"JCEKS");
        keyStore.load(keystoreIn, keystorePW.toCharArray());   //Keystore password
        KeyStore.PasswordProtection keyPassword =       //Key password
                new KeyStore.PasswordProtection(keystorePW.toCharArray());

        KeyStore.PrivateKeyEntry privateKeyEntry = (KeyStore.PrivateKeyEntry) keyStore.getEntry(keyAlias, keyPassword);

        java.security.cert.Certificate cert = keyStore.getCertificate(keyAlias);
        PublicKey publicKey = cert.getPublicKey();
        PrivateKey privateKey = privateKeyEntry.getPrivateKey();

        return new KeyPair(publicKey, privateKey);
    }

    public static String encrypt(String plainText, PublicKey publicKey) throws Exception {
        Cipher encryptCipher = Cipher.getInstance("RSA");
        encryptCipher.init(Cipher.ENCRYPT_MODE, publicKey);

        byte[] cipherText = encryptCipher.doFinal(plainText.getBytes(UTF_8));

        return Base64.getEncoder().encodeToString(cipherText);
    }

    public static String decrypt(String cipherText, PrivateKey privateKey) throws Exception {
        byte[] bytes = Base64.getDecoder().decode(cipherText);

        Cipher decriptCipher = Cipher.getInstance("RSA");
        decriptCipher.init(Cipher.DECRYPT_MODE, privateKey);

        return new String(decriptCipher.doFinal(bytes), UTF_8);
    }

    public static String sign(String plainText, PrivateKey privateKey) throws Exception {
        Signature privateSignature = Signature.getInstance("SHA256withRSA");
        privateSignature.initSign(privateKey);
        privateSignature.update(plainText.getBytes(UTF_8));

        byte[] signature = privateSignature.sign();

        return Base64.getEncoder().encodeToString(signature);
    }

    public static boolean verify(String plainText, String signature, PublicKey publicKey) throws Exception {
        Signature publicSignature = Signature.getInstance("SHA256withRSA");
        publicSignature.initVerify(publicKey);
        publicSignature.update(plainText.getBytes(UTF_8));

        byte[] signatureBytes = Base64.getDecoder().decode(signature);

        return publicSignature.verify(signatureBytes);
    }

}
