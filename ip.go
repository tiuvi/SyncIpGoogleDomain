package ip

import(
	"time"
	"net"
	"net/http"
	"io"
	."tiuvi/core/dac"
	"strconv"
	"strings"
)

//SpaceFile Permanente
var ipFile *PublicSpaceFile
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
	ipFile = NewSfPermBytes(map[string]int64{"ip":16} , nil , "ip", "ip")

	//Revisamos nuestra ip haciendo una busqueda inversa dns.
	iprecords, err := net.LookupIP(dominio)
	if goError  &&
	ipFile.NRESM(goError , err.Error(),file , folder){}

	//Recorremos las ip y nos quedamos con la ip4
	for _, ip := range iprecords {

		if len(ip.To4()) == 4 {

			ipAfter := ip.String()
			//Escribe esta ip en el campo
			ipFile.SetOneFieldString("ip",ipAfter)
	
		}
	}
	
	//Ejecutamos este bucle cada 20 seugndos
	tikect := time.Tick(20 * time.Second)
	for range tikect {

		//Lee la ip en el campo
		ipAfter := ipFile.GetOneFieldString("ip")
	
		//Revisa tu ip en google
		url := "https://domains.google.com/checkip"
		//Falla mucho aunque digan que no.
		//url := "https://api.ipify.org"

		resp, err := http.Get(url)
		if goError  &&
		ipFile.NRESM(goError , err.Error(),file , folder) ||
		ipFile.NRESM(resp.StatusCode > 299 ,"Response failed with status code:" + strconv.Itoa(resp.StatusCode) ,file , folder){
			continue
		}
		defer resp.Body.Close()



		body, err := io.ReadAll(resp.Body)
		ipNow := string(body)
		if goError  &&
		ipFile.NRESM(goError , err.Error(),file , folder) ||
		ipFile.NRESM(!checkIPAddress(ipNow) ,"Respuesta erronea, no es una ip valida." ,file , folder){
			continue
		}
	
		
	
		if ipNow != ipAfter {

			url := strings.Join([]string{ "https://", user , ":" , pass , 
			"@domains.google.com/nic/update?hostname=" , dominio , "&myip=" , ipNow } ,"")

			resp, err := http.Get(url)
			if goError  &&
			ipFile.NRESM(goError , err.Error(),file , folder) || 
			ipFile.NRESM(resp.StatusCode > 299 ,"Response failed with status code:" + strconv.Itoa(resp.StatusCode) ,file , folder){
				continue
			}
	
			body, err := io.ReadAll(resp.Body)
			if goError  &&
			ipFile.NRESM(goError , err.Error(),file , folder){
				continue
			}
			defer resp.Body.Close()

			ipFile.SetOneFieldString("ip",ipNow)

			ipFile.NRESM(true ,"IpUpdate: " + string(body),file , folder)
		}
	}
}

func checkIPAddress(ip string)bool {

	if net.ParseIP(ip) != nil {
		return true
	} 

	return false
}