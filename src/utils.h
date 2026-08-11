#ifndef UTILS_H
#define UTILS_H

int generate_api_key(char key[64]);
int hash_api_key(uint8_t hash[32], const char key[64]);

void hex_encode(char* dst, const uint8_t* src, size_t src_len);
int  hex_decode(uint8_t* dst, const char* src, size_t src_len);

#endif
