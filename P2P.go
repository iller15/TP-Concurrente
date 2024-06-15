package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

var addrs []string
var hostaddr string
var workArrayS []string
var workArrayF []float64

const (
	registerport  = 8000
	notifyport    = 8001
	arraysendport = 8002
)

func main() {
	//descubrir la IP
	hostaddr = descubrirIP()
	hostaddr = strings.TrimSpace(hostaddr)
	fmt.Println("Mi IP: ", hostaddr)
	//Modo Escuchar
	go registerServer() //servicio de enrrollamiento de un nuevo nodo

	//modo cliente
	//menu para conexión
	br := bufio.NewReader(os.Stdin) //creamos buffer
	fmt.Print("Ingrese la IP del nodo remoto: ")
	remoteIP, _ := br.ReadString('\n')
	remoteIP = strings.TrimSpace(remoteIP)

	if remoteIP != "" {
		registerClient(remoteIP)
	}

	//servicio de notificación modo escucha
	notifyServer()

}

func notifyServer() {
	hostname := fmt.Sprintf("%s:%d", hostaddr, notifyport)
	ls, errLisNotServer := net.Listen("tcp", hostname)
	if errLisNotServer != nil {
		fmt.Println("Error listen Notify Server: ", errLisNotServer)
	}
	defer ls.Close()
	for {
		conn, _ := ls.Accept()
		go handleNotify(conn)
	}

}

func handleArray(conn net.Conn) {
	defer conn.Close()
	br := bufio.NewReader(conn)
	array_strin, _ := br.ReadString('\n')
	workArrayS = append(workArrayS, array_strin)
}

func handleNotify(conn net.Conn) {
	defer conn.Close()
	br := bufio.NewReader(conn)
	ip, _ := br.ReadString('\n') // leemos todo hast jump line en esta caso la IP del nuevo nodo
	ip = strings.TrimSpace(ip)

	//actualizar bitacora de ips con la nueva
	addrs = append(addrs, ip)
	fmt.Println(addrs)
}

func registerClient(remoteIP string) {
	remotehost := fmt.Sprintf("%s:%d", remoteIP, registerport)
	println(remotehost)
	//conectarme al nodo remoto
	conn, errDialClient := net.Dial("tpc", remotehost)
	if errDialClient != nil {
		fmt.Print("Error Dial Cliente: ", errDialClient)
	}
	defer conn.Close()
	fmt.Fprintf(conn, "%s\n", hostaddr) //envio la ip al nodo remota para el enrolamiento

	//recibir la bitácora de ips y actualizarla
	br := bufio.NewReader(conn)
	bitacoraIPs, _ := br.ReadString('\n') //la bitacora de IPs
	var bitaIps []string
	json.Unmarshal([]byte(bitacoraIPs), &bitaIps) //deserializar nuestra bitacora y lo ponemos en bitaIps || la bitacora la recibimos del nodo servidor como respuesta a registarnos
	//actualizar bitacora
	addrs = append(bitaIps, remoteIP) //addicionar en el arreglos ||no hay problema si un nodo se registra una sola vez por red.
	println(addrs)

}

func registerServer() {
	hostname := fmt.Sprintf("%s:%d", hostaddr, registerport)
	ls, errLs := net.Listen("tcp", hostname)
	if errLs != nil {
		fmt.Println("Error Listen registerServer(): ", errLs)
	}
	defer ls.Close()

	//manejar conexiones entrantes
	for {
		conn, _ := ls.Accept()
		go handleRegister(conn) //concurrente para soportar alto volumen
	}

}

func handleRegister(conn net.Conn) {
	defer conn.Close()
	br := bufio.NewReader(conn)
	ip, _ := br.ReadString('\n') // leemos todo hast jump line en esta caso la IP del nuevo nodo
	ip = strings.TrimSpace(ip)

	//enviar al nuevo nodo, la bitacora actual de IPs
	bytesjson, errJson := json.Marshal(addrs)

	if errJson != nil {
		fmt.Println("Error Json: ", errJson)
	}

	fmt.Fprintf(conn, "%s\n", string(bytesjson)) //enviamos nuestra bitacora serializada en string a el nodo que se estra regsitrando

	//notificar al resto de nodos el registro del nuevo nodo
	notifyAll(ip)

	//actualizar bitacora
	addrs = append(addrs, ip)

	println(addrs)
}

func notifyAll(ip string) {
	for _, addres := range addrs {
		notify(addres, ip)
	}
}

func notify(addr, ip string) {
	remotehost := fmt.Sprintf("%s:%d", addr, notifyport)
	conn, errDial := net.Dial("tcp", remotehost)
	if errDial != nil {
		println("Error Dial notify: ", errDial, remotehost)
		print(addrs)
	}
	defer conn.Close()

	//enviamos la ip nueva al addrs
	fmt.Fprintln(conn, ip) //al enviar asgurarse de tener un salto de linea al final
}

func descubrirIP() string {
	//interfas de red
	ifaces, _ := net.Interfaces()
	for _, i := range ifaces { //interfaces de red
		//if strings.HasPrefix(i.Name, "Ethernet") {
		//if strings.HasPrefix(i.Name, "eth") { USAR ESTE SI ES VM
		if strings.Contains(i.Name, "Wireless LAN adapter Wi-Fi") { //USAR ESTE PARA CUANDO USAS WIFI
			//solo si tiene "Ethernet" al inicio
			addrs, _ := i.Addrs()
			for _, addr := range addrs {
				switch t := addr.(type) {
				case *net.IPNet: //si el tipo es ip? o es en caso de emergencia
					if t.IP.To4() != nil {
						return t.IP.To4().String() //retornamos la ip ethernet version 4
					} //if el tipo de ip es ipv4
				}
			}
		}
	}
	return "127.0.0.1" //default ip de un localhost si no está conectado
}
