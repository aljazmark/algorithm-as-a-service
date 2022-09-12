#include "sort.h"
#include <stdio.h>
#include <stdlib.h>
#include <string.h>




void insertionSort(int array[], int n) 
{   
    int i, element, j; 
    for (i = 1; i < n; i++) { 
        element = array[i]; j = i - 1; 
        while (j >= 0 && array[j] > element) { 
            array[j + 1] = array[j]; 
        j = j - 1; 
    } 
    array[j + 1] = element; 
    };  
}
void insertionSortReverse(int array[], int n) 
{   
    int i, element, j; 
    for (i = 1; i < n; i++) { 
        element = array[i]; j = i - 1; 
        while (j >= 0 && array[j] < element) { 
            array[j + 1] = array[j]; 
        j = j - 1; 
    } 
    array[j + 1] = element; 
    };  
}
void swapEl(int *a, int *b) 
{ 
    int temp = *a; 
    *a = *b; 
    *b = temp; 
} 
void bubbleSort(int array[], int n) 
{ 
    int i, j; 
    for (i = 0; i < n-1; i++)  {
        for (j = 0; j < n-i-1; j++){
            if (array[j] > array[j+1]) {
                swapEl(&array[j], &array[j+1]);   
            }         
        }         
    }     
}
void bubbleSortReverse(int array[], int n) 
{ 
    int i, j; 
    for (i = 0; i < n-1; i++)  {
        for (j = 0; j < n-i-1; j++){
            if (array[j] < array[j+1]) {
                swapEl(&array[j], &array[j+1]);   
            }         
        }         
    }     
}

void selectionSort(int array[], int n) 
{ 
int i, j, min_element; 
    for (i = 0; i < (n-1); i++) {
        min_element = i; 
        for (j = i+1; j < n; j++) {
            if (array[min_element] >array[j]) {
                min_element = j;               
            }
            
                   
        }
        if(min_element != i){
                swapEl(&array[min_element], &array[i]);
            }            
    }
}
void selectionSortReverse(int array[], int n) 
{ 
int i, j, min_element; 
    for (i = 0; i < (n-1); i++) {
        min_element = i; 
        for (j = i+1; j < n; j++) {
            if (array[min_element] <array[j]) {
                min_element = j;               
            }
            
                   
        }
        if(min_element != i){
                swapEl(&array[min_element], &array[i]);
            }            
    }
}
int partition (int arr[], int lowIndex, int highIndex) 
{ 
    int pivotElement = arr[highIndex];
    int i = (lowIndex - 1); 
    for (int j = lowIndex; j <= highIndex- 1; j++) { 
        if (arr[j] <= pivotElement) { 
            i++; 
            swapEl(&arr[i], &arr[j]); 
        } 
    } 
    swapEl(&arr[i + 1], &arr[highIndex]); 
    return (i + 1); 
} 
void quickSort(int arr[], int lowIndex, int highIndex) 
{ 
    if (lowIndex < highIndex) { 
        int pivot = partition(arr, lowIndex, highIndex); 
        quickSort(arr, lowIndex, pivot - 1); 
        quickSort(arr, pivot + 1, highIndex); 
    } 
}
int partitionReverse (int arr[], int lowIndex, int highIndex) 
{ 
    int pivotElement = arr[highIndex];
    int i = (lowIndex - 1); 
    for (int j = lowIndex; j <= highIndex- 1; j++) { 
        if (arr[j] > pivotElement) { 
            i++; 
            swapEl(&arr[i], &arr[j]); 
        } 
    } 
    swapEl(&arr[i + 1], &arr[highIndex]); 
    return (i + 1); 
} 
void quickSortReverse(int arr[], int lowIndex, int highIndex) 
{ 
    if (lowIndex < highIndex) { 
        int pivot = partitionReverse(arr, lowIndex, highIndex); 
        quickSortReverse(arr, lowIndex, pivot - 1); 
        quickSortReverse(arr, pivot + 1, highIndex); 
    } 
}
void merge(int arr[], int l, int m, int r) 
{ 
    int i, j, k; 
    int n1 = m - l + 1; 
    int n2 =  r - m; 
    int L[n1], R[n2]; 
    for (i = 0; i < n1; i++) {
        L[i] = arr[l + i]; 
    }   
    for (j = 0; j < n2; j++) {
        R[j] = arr[m + 1+ j]; 
    }   
    i = 0; 
    j = 0; 
    k = l; 
    while (i < n1 && j < n2) 
    { 
        if (L[i] <= R[j]) 
        { 
            arr[k] = L[i]; 
            i++; 
        } 
        else
        { 
            arr[k] = R[j]; 
            j++; 
        } 
        k++; 
    } 
    while (i < n1) 
    { 
        arr[k] = L[i]; 
        i++; 
        k++; 
    } 
    while (j < n2) 
    { 
        arr[k] = R[j]; 
        j++; 
        k++; 
    } 
}
void mergeSort(int arr[], int l, int r) 
{ 
    if (l < r) 
    { 
    int m = l+(r-l)/2; 
    mergeSort(arr, l, m); 
    mergeSort(arr, m+1, r); 
    merge(arr, l, m, r); 
    }
}
void mergeReverse(int arr[], int l, int m, int r) 
{ 
    int i, j, k; 
    int n1 = m - l + 1; 
    int n2 =  r - m; 
    int L[n1], R[n2]; 
    for (i = 0; i < n1; i++) {
        L[i] = arr[l + i]; 
    }   
    for (j = 0; j < n2; j++) {
        R[j] = arr[m + 1+ j]; 
    }   
    i = 0; 
    j = 0; 
    k = l; 
    while (i < n1 && j < n2) 
    { 
        if (L[i] > R[j]) 
        { 
            arr[k] = L[i]; 
            i++; 
        } 
        else
        { 
            arr[k] = R[j]; 
            j++; 
        } 
        k++; 
    } 
    while (i < n1) 
    { 
        arr[k] = L[i]; 
        i++; 
        k++; 
    } 
    while (j < n2) 
    { 
        arr[k] = R[j]; 
        j++; 
        k++; 
    } 
}
void mergeSortReverse(int arr[], int l, int r) 
{ 
    if (l < r) 
    { 
    int m = l+(r-l)/2; 
    mergeSortReverse(arr, l, m); 
    mergeSortReverse(arr, m+1, r); 
    mergeReverse(arr, l, m, r); 
    }
}

