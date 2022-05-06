package ip

import (
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
	. "tiuvi/core/dac"
)

//SpaceFile Permanente
var IpFile *PublicSpaceFile
var file , folder = "ip" , "ipErrors"

//Errores
var err error
var goError = err != nil	

//Inicia un dac basico
func NewDacForIp(path string ){

	//Crea dac en esta ruta
	NewBasicDac(path)
}

//Funcion que actualiza la ip en google domain cada 20 segundos
func InitUpdateIp(user string, pass string , dominio string){
	
	//Crea un espacio permanente con un campo ip de 16 bytes en dac/ip/ip.dacbyte
	IpFile = NewSfPermBytes(map[string]int64{"ip":16} , nil , "ip", "ip")

	//Revisamos nuestra ip haciendo una busqueda inversa dns.
	iprecords, err := net.LookupIP(dominio)
	if goError  &&
	IpFile.NRESM(goError , err.Error(),file , folder){}

	//Recorremos las ip y nos quedamos con la ip4
	for _, ip := range iprecords {

		if len(ip.To4()) == 4 {

			ipAfter := ip.String()
			//Escribe esta ip en el campo
			IpFile.SetOneFieldString("ip",ipAfter)
	
		}
	}

	 go UpdateIp(user , pass, dominio)

}

func UpdateIp(user string, pass string , dominio string){

	var ipAfter string
	var url string = "https://domains.google.com/checkip"
	var respIp *http.Response
	var respUpdIp *http.Response
	//Ejecutamos este bucle cada 20 seugndos
	tikect := time.Tick(20 * time.Second)
	for range tikect {

		//Lee la ip en el campo
		ipAfter = IpFile.GetOneFieldString("ip")
	

		respIp, err = http.Get(url)
		if goError  &&
		IpFile.NRESM(goError , err.Error() ,file , folder) ||
		IpFile.NRESM(respIp.StatusCode > 299 ,"Response failed with status code:" + strconv.Itoa(respIp.StatusCode) ,file , folder){
			
			continue
		}
			



		body, err := io.ReadAll(respIp.Body)
		ipNow := string(body)
		if goError  &&
		IpFile.NRESM(goError , err.Error(),file , folder) ||
		IpFile.NRESM(!checkIPAddress(ipNow) ,"Respuesta erronea, no es una ip valida." ,file , folder){
			continue
		}
	
		err = respIp.Body.Close()
		if goError  &&
		IpFile.NRESM(goError , err.Error(),file , folder){
			continue
		}
	
	
		
		if ipNow != ipAfter {

			url := strings.Join([]string{ "https://", user , ":" , pass , 
			"@domains.google.com/nic/update?hostname=" , dominio , "&myip=" , ipNow } ,"")

			respUpdIp, err = http.Get(url)	 
			if goError  &&
			IpFile.NRESM(goError , err.Error() ,file , folder) || 
			IpFile.NRESM(respUpdIp.StatusCode > 299 ,"Response failed with status code:" + strconv.Itoa(respUpdIp.StatusCode) ,file , folder){
				
				continue
			}
	
			body, err := io.ReadAll(respUpdIp.Body)
			if goError  &&
			IpFile.NRESM(goError , err.Error(),file , folder){
				continue
			}
			
			err = respUpdIp.Body.Close()
			if goError  &&
			IpFile.NRESM(goError , err.Error(),file , folder){
				continue
			}


			IpFile.SetOneFieldString("ip",ipNow)

			IpFile.NRESM(true ,"IpUpdate: " + string(body),file , folder)

		}
	}
}

func checkIPAddress(ip string)bool {

	if net.ParseIP(ip) != nil {
		return true
	} 

	return false
}