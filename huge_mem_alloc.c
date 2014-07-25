#include<stdlib.h>
#include<unistd.h>
#include<string.h>
#include<stdio.h>

void main() {
        size_t mem_size = sizeof(int) * 1024 * 1024 * 30;
            int * p = (int *)malloc(mem_size);
            if (p  ==  NULL) {
                        printf("mem alloc fail\n");

            }else{
                        printf("mem alloc size:%d\n", mem_size);

            }
                bzero(p, mem_size);
                    sleep(1000);

}
