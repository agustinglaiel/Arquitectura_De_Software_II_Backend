package utils

import (
	"time"

	"log"

	"github.com/allegro/bigcache"
)

var cache *bigcache.BigCache

func InitCache(cacheDuration time.Duration) {
    var err error
    cache, err = bigcache.NewBigCache(bigcache.Config{
        Shards:             1024,           // Número de shards (particiones), aumenta para mejorar la concurrencia
        LifeWindow:         cacheDuration,  // Tiempo de vida de cada clave en la caché
        CleanWindow:        5 * time.Minute, // Intervalo de limpieza de claves expiradas
        MaxEntriesInWindow: 1000 * 10 * 60,  // Máximo número de entradas en la ventana de vida
        MaxEntrySize:       500,            // Máximo tamaño en bytes de cada entrada
        Verbose:            true,           // Modo verbose para mostrar información adicional
    })

    if err != nil {
        log.Fatalf("Failed to initialize cache: %v", err)
    }
}

func Set(key string, entry []byte) error {
    return cache.Set(key, entry)
}

func Get(key string) ([]byte, error) {
    return cache.Get(key)
}

func Delete(key string) error {
    return cache.Delete(key)
}

func Reset() {
    cache.Reset() // Limpia toda la caché
}
