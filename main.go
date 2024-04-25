package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const fileName = "larger.txt" // Nombre del archivo de texto a analizar

func main() {
	// Calcular el tiempo de inicio de ejecución
	startTime := time.Now()
	totalWordHashMap := make(map[string]int) // Mapa para almacenar el recuento total de palabras

	numWorkers := runtime.NumCPU()                        // Número de núcleos de CPU disponibles en el sistema
	var wg sync.WaitGroup                                 // WaitGroup para sincronizar las goroutines
	interimResultFromWorkers := make(chan map[string]int) // Canal para recibir resultados parciales de las goroutines

	fileInfo, err := os.Stat(fileName) // Obtener información sobre el archivo
	if err != nil {
		panic("Error getting file info")
	}
	fileSize := fileInfo.Size()               // Tamaño del archivo en bytes
	chunkSize := fileSize / int64(numWorkers) // Tamaño de los fragmentos del archivo para cada trabajador

	fmt.Printf("\nfile info: size %v\n", fileSize) // Imprimir información sobre el archivo

	// Iniciar goroutines para procesar diferentes partes del archivo
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		start := int64(i) * chunkSize
		var end int64
		if i == numWorkers-1 {
			end = fileSize
		} else {
			end = start + chunkSize
		}
		go readFileFromCertainChunk(fileName, start, end, interimResultFromWorkers, &wg)
	}

	// Crear una goroutine para cerrar el canal cuando todos los trabajadores hayan terminado
	go func() {
		wg.Wait()
		// Tiempo transcurrido
		elapsedTime := time.Since(startTime).Milliseconds()
		println("Time elapsed in milliseconds: ", elapsedTime)
		close(interimResultFromWorkers) // Cerrar el canal
	}()

	// Recibir resultados parciales de las goroutines y agregarlos al mapa totalWordHashMap
	for result := range interimResultFromWorkers {
		for key, value := range result {
			totalWordHashMap[key] += value
		}
	}

	WriteResultsToFile(totalWordHashMap) // Escribir resultados en un archivo
}

// Función para leer una parte específica del archivo y contar la ocurrencia de palabras
func readFileFromCertainChunk(fileName string, start, end int64, result chan<- map[string]int, wg *sync.WaitGroup) {
	defer wg.Done() // Declarar que la goroutine ha terminado

	file, err := os.Open(fileName) // Abrir el archivo
	if err != nil {
		panic("Error opening file")
	}
	defer file.Close()

	file.Seek(start, 0)               // Ir a la posición de inicio en el archivo
	buffer := make([]byte, end-start) // Crear un búfer para almacenar los datos del archivo
	_, err = file.Read(buffer)        // Leer datos del archivo
	
	if err != nil {
		panic("Error reading file")
	}

	words := string(buffer) // Convertir los datos del archivo en una cadena
	fruits := "manzanas bananas naranjas peras uvas kiwis sandias melones mangos fresas frambuesas papayas limones cerezas moras aguacates cocos"
	fruitList := strings.Fields(fruits) // Convertir la cadena de frutas en una lista de palabras

	wordCount := make(map[string]int) // Mapa para almacenar el recuento de palabras
	for _, fruit := range fruitList {
		wordCount[fruit] = countWord(words, fruit) // Contar la ocurrencia de cada palabra en la parte del archivo
	}

	result <- wordCount // Enviar el recuento de palabras al canal de resultados
}

// Función para contar la ocurrencia de una palabra en una cadena
func countWord(words, word string) int {
	return strings.Count(words, word)
}

// Función para escribir los resultados en un archivo
func WriteResultsToFile(resultMap map[string]int) {
	file, err := os.Create("resultados.txt") // Crear un archivo para escribir resultados
	if err != nil {
		panic("Error creating file")
	}
	defer file.Close()

	// Escribir los resultados en el archivo
	for fruit, count := range resultMap {
		_, err := fmt.Fprintf(file, "%s: %d\n", fruit, count)
		if err != nil {
			panic("Error writing to file")
		}
	}
}
