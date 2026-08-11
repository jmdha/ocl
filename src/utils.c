#include <stdint.h>
#include <stddef.h>
#include <openssl/rand.h>
#include <string.h>

#include "utils.h"

int generate_api_key(char key[64]) {
	unsigned char raw[31];
	static const char hex_chars[] = "0123456789abcdef";
	
	if (RAND_bytes(raw, sizeof(raw)) != 1)
		return -1;
	
	for (size_t i = 0; i < sizeof(raw); i++) {
		key[i * 2]     = hex_chars[(raw[i] >> 4) & 0x0F];
		key[i * 2 + 1] = hex_chars[raw[i] & 0x0F];
	}
	
	key[sizeof(raw) * 2] = '\0'; /* key[62] = '\0' */
	
	return 0;
}

int hash_api_key(uint8_t hash[32], const char key[64]) {
	unsigned int  len = 0;
	EVP_MD_CTX   *ctx = EVP_MD_CTX_new();
	
	if (!ctx) return -1;
	
	int ok = EVP_DigestInit_ex(ctx, EVP_sha256(), NULL) &&
	         EVP_DigestUpdate(ctx, key, strlen(key)) &&
	         EVP_DigestFinal_ex(ctx, hash, &len);
	
	EVP_MD_CTX_free(ctx);
	
	if (!ok || len != 32) return -1;
	
	return 0;
}

void hex_encode(char* dst, const uint8_t* src, size_t src_len) {
	static const char lut[] = "0123456789abcdef";
	for (size_t i = 0; i < src_len; i++) {
		dst[i * 2]     = lut[src[i] >> 4];
		dst[i * 2 + 1] = lut[src[i] & 0x0f];
	}
	dst[src_len * 2] = '\0';
}

int hex_decode(uint8_t* dst, const char* src, size_t src_len) {
	if (src_len % 2 != 0) return -1;
	for (size_t i = 0; i < src_len / 2; i++) {
		int hi = -1, lo = -1;
		char c;

		c = src[i * 2];
		if (c >= '0' && c <= '9') hi = c - '0';
		else if (c >= 'a' && c <= 'f') hi = c - 'a' + 10;
		else if (c >= 'A' && c <= 'F') hi = c - 'A' + 10;

		c = src[i * 2 + 1];
		if (c >= '0' && c <= '9') lo = c - '0';
		else if (c >= 'a' && c <= 'f') lo = c - 'a' + 10;
		else if (c >= 'A' && c <= 'F') lo = c - 'A' + 10;

		if (hi < 0 || lo < 0) return -1;
		dst[i] = (uint8_t)((hi << 4) | lo);
	}
	return 0;
}
