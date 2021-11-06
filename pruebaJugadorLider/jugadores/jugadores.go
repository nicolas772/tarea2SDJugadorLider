package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	pb "example.com/go-msgs-grpc/msgs"
	"google.golang.org/grpc"
)

const (
	address = "10.6.40.193:60051"
)

func printmenuInicial() string {
	menu :=
		`
		Bienvenido al juego del calamar
		[ 1 ] Enviar peticion para jugar
		[ otro ] Salir
		Selecciona una opcion
	`
	fmt.Print(menu)
	reader := bufio.NewReader(os.Stdin)

	entrada, _ := reader.ReadString('\n')          // Leer hasta el separador de salto de línea
	eleccion := strings.TrimRight(entrada, "\r\n") // Remover el salto de línea de la entrada del usuario
	return eleccion
}

func Instrucciones1() {
	menu :=
		`
	El juego consiste en que cada jugador deberá elegir números entre el 1 y el 10 en diferentes rondas hasta sumar
	21. Además de esto, en cada ronda el Lider elegirá al azar un número entre el 6 y el 10.
	En cada ronda se verificará que jugadores eligieron un número igual o mayor al del Líder. Quienes lo hayan
	hecho quedan eliminados del juego.
	El juego durará máximo 4 rondas. Todos los jugadores que no logren llegar a 21 antes de las 4 rondas
	serán eliminados.
	`
	fmt.Print(menu)
}

func eleccionUsuario1() string {
	fmt.Println("--------------------------------")
	fmt.Println("Escribe un numero entre 1 y 10")
	reader := bufio.NewReader(os.Stdin)

	entrada, _ := reader.ReadString('\n')          // Leer hasta el separador de salto de línea
	eleccion := strings.TrimRight(entrada, "\r\n") // Remover el salto de línea de la entrada del usuario
	return eleccion
}

func main() {

	eleccion_inicial := printmenuInicial()

	if eleccion_inicial == "1" {
		fmt.Println("Enviando solicitud para jugar")
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewLiderManagementClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.PeticionJugar(ctx, &pb.Peticion{NombreJugador: "Nicolas"})
		if err != nil {
			log.Fatalf("no se pudo enviar la peticion para jugar: %v", err)
		}
		if r.GetPuedeJugar() == 0 {
			log.Printf("No puedes jugar, los cupos estan completos")
		} else {
			log.Printf("Bienvenido a la primera Etapa: Luz Roja, Luz Verde")
			log.Printf("Tu id asignado es: %v", r.GetIdAsignado())
			my_id := r.GetIdAsignado()
			Instrucciones1()
			flag := 1

			for flag == 1 {
				aux := eleccionUsuario1()
				eleccion_lvlr, _ := strconv.Atoi(aux)
				ctx1, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()
				r1, err := c.Enviarjugada(ctx1, &pb.Jugada{IdJugador: my_id, Contenido: int32(eleccion_lvlr)})
				if err != nil {
					log.Fatalf("no se pudo enviar la jugada: %v", err)
				}
				sigo_jugando := r1.GetSiguejugando()
				flag = int(sigo_jugando)
			}

		}
	} else {
		fmt.Println("Hasta pronto!")
	}
	//--------------Intento para 16 jugadores ---------------------------------------------------

	/*jugador := 1
	for jugador <= 15 {
		go func() {
			aux1 := 6000 + jugador
			address1 := "localhost:" + strconv.Itoa(aux1)
			conn, err := grpc.Dial(address1, grpc.WithInsecure(), grpc.WithBlock())
			if err != nil {
				log.Fatalf("did not connect: %v", err)
			}
			defer conn.Close()
			c := pb.NewLiderManagementClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			r, err := c.PeticionJugar(ctx, &pb.Peticion{NombreJugador: "Jugador " + strconv.Itoa(jugador)})
			if err != nil {
				log.Fatalf("no se pudo enviar la peticion para jugar: %v", err)
			}
			if r.GetPuedeJugar() == 0 {
				log.Printf("No puedes jugar, los cupos estan completos")
			} else {
				//log.Printf("Bienvenido a la primera Etapa: Luz Roja, Luz Verde")
				//log.Printf("Tu id asignado es: %v", r.GetIdAsignado())
				my_id := r.GetIdAsignado()
				flag := 1
				for flag == 1 {
					rand.Seed(int64(time.Now().UnixNano()))
					aux := rand.Intn(10)
					jugada_bot := aux + 1
					ctx2, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()
					r1, err := c.Enviarjugada(ctx2, &pb.Jugada{IdJugador: my_id, Contenido: int32(jugada_bot)})
					if err != nil {
						log.Fatalf("no se pudo enviar la jugada: %v", err)
					}
					sigo_jugando1 := r1.GetSiguejugando()
					flag = int(sigo_jugando1)
				}

			}
		}()
		jugador++
	}*/
}
