package com.fjn.magazinereturncandidate.common.utils;

import android.util.Base64;

import java.math.BigInteger;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.security.SecureRandom;
import java.util.Arrays;

import javax.crypto.Cipher;
import javax.crypto.spec.IvParameterSpec;
import javax.crypto.spec.SecretKeySpec;

/**
 * Common encrypt info user and decypt info user
 * Created by cong-pv on 2018/10/16.
 */

public class EnDecryptInfoCommon {


    /**
     * Function  EncryptMD5
     *
     * @param strEncrypt string encrypt (password....)
     * @return String is encrypt success
     */
    public static String encryptMD5(String strEncrypt) {
        try {
            MessageDigest md = MessageDigest.getInstance("MD5");
            byte[] messageDigest = md.digest(strEncrypt.getBytes());
            BigInteger number = new BigInteger(1, messageDigest);
            String hashtext = number.toString(16);
            while (hashtext.length() < 32) {
                hashtext = "0" + hashtext;
            }
            return hashtext;
        } catch (NoSuchAlgorithmException e) {
            e.printStackTrace();
            return null;
        }
    }

    /**
     * Decrypt info user
     * @param strEnDeCrypt string decrypt (user_id, license...)
     * @return String is decrypt success
     */
    public static String decryptString(String strEnDeCrypt) {
        String TOKEN_KEY = "fqJfdzGDvfwbedsKSUGty3VZ9taXxMVw";
        try {
            byte[] ivAndCipherText = Base64.decode(strEnDeCrypt, Base64.NO_WRAP);
            byte[] iv = Arrays.copyOfRange(ivAndCipherText, 0, 16);
            byte[] cipherText = Arrays.copyOfRange(ivAndCipherText, 16, ivAndCipherText.length);

            Cipher cipher = Cipher.getInstance("AES/CBC/PKCS5Padding");
            cipher.init(Cipher.DECRYPT_MODE, new SecretKeySpec(TOKEN_KEY.getBytes("utf-8"), "AES"), new IvParameterSpec(iv));
            return new String(cipher.doFinal(cipherText), "utf-8");
        } catch (Exception e) {
            e.printStackTrace();
            return null;
        }
    }

    /**
     * Function encrypt DES
     * @param strEnDeCrypt string encrypt DES
     * @return string encrypt success
     */
    //encrypt and decrypt DES
    public String encryptString(String strEnDeCrypt) {

        String TOKEN_KEY = "fqJfdzGDvfwbedsKSUGty3VZ9taXxMVw";
        try {
            byte[] iv = new byte[16];
            new SecureRandom().nextBytes(iv);
            Cipher cipher = Cipher.getInstance("AES/CBC/PKCS5Padding");
            cipher.init(Cipher.ENCRYPT_MODE, new SecretKeySpec(TOKEN_KEY.getBytes("utf-8"), "AES"), new IvParameterSpec(iv));
            byte[] cipherText = cipher.doFinal(strEnDeCrypt.getBytes("utf-8"));
            byte[] ivAndCipherText = getCombinedArray(iv, cipherText);
            return Base64.encodeToString(ivAndCipherText, Base64.NO_WRAP);
        } catch (Exception e) {
            e.printStackTrace();
            return null;
        }
    }

    private static byte[] getCombinedArray(byte[] one, byte[] two) {

        byte[] combined = new byte[one.length + two.length];
        for (int i = 0; i < combined.length; ++i) {
            combined[i] = i < one.length ? one[i] : two[i - one.length];
        }
        return combined;
    }


}
