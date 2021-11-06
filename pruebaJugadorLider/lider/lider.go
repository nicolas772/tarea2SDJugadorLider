package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"time"

	pb "example.com/go-msgs-grpc/msgs"
	"google.golang.org/grpc"
)

const (
	port                 = ":60051"
	cantidadMaxJugadores = 4
)

func NewLiderManagementServer() *LiderManagementServer {
	return &LiderManagementServer{
		cantidad_jugadores_actual: 0,
		puntaje_jugador1:          0,
		ronda:                     0,
	}
}

type LiderManagementServer struct {
	pb.UnimplementedLiderManagementServer
	cantidad_jugadores_actual int
	puntaje_jugador1          int
	ronda                     int
}

func (server *LiderManagementServer) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterLiderManagementServer(s, &LiderManagementServer{})
	log.Printf("server listening at %v", lis.Addr())
	return s.Serve(lis)

}

func (s *LiderManagementServer) Enviarjugada(ctx context.Context, in *pb.Jugada) (*pb.Respuestalider, error) {
	rand.Seed(int64(time.Now().UnixNano()))
	id_jugador := in.GetIdJugador()
	jugada := in.GetContenido()
	log.Printf("jugada jugador %v : %v", id_jugador, jugada)
	s.puntaje_jugador1 += int(jugada)
	s.ronda += 1
	log.Printf("RONDA NUMERO: %d", s.ronda)
	aux := rand.Intn(5)
	jugada_lider := aux + 6
	log.Printf("jugada lider: %d", jugada_lider)
	if int32(jugada_lider) <= jugada {
		log.Printf("jugada lider: %d", jugada_lider)
		log.Printf("Jugador %d eliminado", id_jugador)
		return &pb.Respuestalider{Siguejugando: 0, Pasaronda: 0}, nil
	} else {
		if s.ronda < 4 {
			if s.puntaje_jugador1 < 21 {
				log.Printf("puntaje actual jugador 1: %d", s.puntaje_jugador1)
				log.Printf("jugador 1 debe seguir jugando")
				return &pb.Respuestalider{Siguejugando: 1, Pasaronda: 0}, nil
			} else {
				log.Printf("jugador 1 supero la prueba. FELICITACIONES!")
				return &pb.Respuestalider{Siguejugando: 0, Pasaronda: 1}, nil
			}
		}
		if s.ronda == 4 {
			if s.puntaje_jugador1 < 21 {
				log.Printf("Jugador eliminado, ya no quedan mas rondas")
				return &pb.Respuestalider{Siguejugando: 0, Pasaronda: 0}, nil
			} else {
				log.Printf("jugador 1 supero la prueba. FELICITACIONES!")
				return &pb.Respuestalider{Siguejugando: 0, Pasaronda: 1}, nil
			}
		}
		if s.ronda > 4 {
			log.Printf("Jugador eliminado, ya no quedan mas rondas")
			return &pb.Respuestalider{Siguejugando: 0, Pasaronda: 0}, nil
		}
	}
	return &pb.Respuestalider{Siguejugando: 0, Pasaronda: 0}, nil
}

func (s *LiderManagementServer) PeticionJugar(ctx context.Context, in *pb.Peticion) (*pb.RespuestaPeticion, error) {
	if s.cantidad_jugadores_actual <= cantidadMaxJugadores {
		s.cantidad_jugadores_actual += 1
		log.Printf("Bienvenido %v", in.GetNombreJugador())
		log.Printf("tu id es: %d", s.cantidad_jugadores_actual)
		return &pb.RespuestaPeticion{PuedeJugar: 1, IdAsignado: int32(s.cantidad_jugadores_actual)}, nil
	} else {
		log.Printf("Lo siento, se han llenado los cupos")
		return &pb.RespuestaPeticion{PuedeJugar: 0, IdAsignado: 0}, nil
	}
}

func main() {

	var lider_mgmt_server *LiderManagementServer = NewLiderManagementServer()
	if err := lider_mgmt_server.Run(); err != nil {
		log.Fatalf("failed to server: %v", err)
	}

	//----------------------Intento para otros 15 jugadores-------------
	/*i := 1
	for i <= 15 {
		go func() {
			var lider_mgmt_server *LiderManagementServer = NewLiderManagementServer()
			aux1 := 6000 + i
			puerto1 := strconv.Itoa(aux1)
			if err := lider_mgmt_server.Run(puerto1); err != nil {
				log.Fatalf("failed to server: %v", err)
			}
		}()
		i++
	}*/

}
