package ip

import(
	"time"
	"net"
	"log"
	"net/http"
	"io"
	."tiuvi/core/dac"
	"strconv"
	"strings"
)

var ipFile *PublicSpaceFile
var file , folder = "ip" , "ipErrors"

var err error
var goError = err != nil	


func NewDacForIp(path string ){

	//Crea dac en esta ruta
	NewBasicDac(path)
}

func InitUpdateIp(user string, pass string , dominio string){
	
	//Crea un espacio permanente con un campo ip de 16 bytes en dac/ip/ip.dacbyte
	ipFile = NewSfPermBytes(map[string]int64{"ip":16} , nil , "ip", "ip")

	iprecords, err := net.LookupIP(dominio)
	if goError  &&
	ipFile.NRESM(goError , err.Error(),file , folder){}

	for _, ip := range iprecords {

		if len(ip.To4()) == 4 {

			ipAfter := ip.String()
			//Escribe esta ip en el campo
			ipFile.SetOneFieldString("ip",ipAfter)
			log.Println(ipAfter)
		}
	}
	
	tikect := time.Tick(20 * time.Second)

	for range tikect {

		//Lee la ip en el campo
		ipAfter := ipFile.GetOneFieldString("ip") 
		log.Println(ipAfter)
		//Alternativa
		//https://domains.google.com/checkip
		url := "https://api.ipify.org"
		resp, err := http.Get(url)
		if goError  &&
		ipFile.NRESM(goError , err.Error(),file , folder){
			continue
		}

		if goError  &&
		ipFile.NRESM(resp.StatusCode > 299 ,"Response failed with status code:" + strconv.Itoa(resp.StatusCode) ,file , folder){
			continue
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if goError  &&
		ipFile.NRESM(goError , err.Error(),file , folder){
			continue
		}

		ipNow := string(body)

		log.Println("ipnow:   " , ipNow ," -ipNowLen: " , len(ipNow))
		log.Println("ipAfter: " ,ipAfter ," -ipNowLen: ", len(ipNow))
		log.Println("ipnow:   " , []byte(ipNow) ," -ipNowLen: " , len(ipNow))
		log.Println("ipAfter: " ,[]byte(ipAfter) ," -ipNowLen: ", len(ipNow))
		log.Println(ipNow != ipAfter )

		if ipNow != ipAfter {

			url := strings.Join([]string{ "https://", user , ":" , pass , 
			"@domains.google.com/nic/update?hostname=" , dominio , "&myip=" , ipNow } ,"")

			resp, err := http.Get(url)
			if goError  &&
			ipFile.NRESM(goError , err.Error(),file , folder){
				continue
			}
	
			body, err := io.ReadAll(resp.Body)
			defer resp.Body.Close()
			ipFile.SetOneFieldString("ip",ipNow)


			log.Println(string(body))

		}

	}
}

