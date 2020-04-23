package api

import (
	"context"
	"github.com/BASChain/go-bas-mail-server/app/cmdcommon"
	"github.com/BASChain/go-bas-mail-server/app/cmdpb"
)

type CmdStringOPSrv struct {
}

func (cso *CmdStringOPSrv) StringOpDo(cxt context.Context, so *cmdpb.StringOP) (*cmdpb.DefaultResp, error) {
	msg := ""
	switch so.Op {
	case cmdcommon.CMD_EXAMPLE:
		msg = showExample(so.Param)
	//case cmdcommon.CMD_DOMAIN:
	//	msg = GetRecords(so.Param)
	//case cmdcommon.CMD_DEAL:
	//	msg = GetDeal(so.Param)
	//case cmdcommon.CMD_ORDER:
	//	msg = GetOrder(so.Param)
	default:
		return encapResp("Command Not Found"), nil
	}

	return encapResp(msg), nil
}

func showExample(example string) string {
	msg := ""

	if example == "" {
		msg += "no example"
	} else {
		msg += example
	}

	//if wallet == "" {
	//	for k, ass := range DataSync.Assets {
	//		msg += "Wallet: " + k.String() + "\r\n"
	//		msg += getAssetInfo(ass)
	//		msg += "\r\n"
	//	}
	//
	//	if msg == "" {
	//		msg = "No assets"
	//	}
	//
	//	return msg
	//}
	//msg += "Wallet: " + wallet + "\r\n"
	//addr := common.HexToAddress(wallet)
	//if a, ok := DataSync.Assets[addr]; !ok {
	//	msg = "NotFound"
	//} else {
	//	msg = getAssetInfo(a)
	//}

	return msg
}

//
//func getAssetInfo(domains []Bas_Ethereum.Hash) string {
//	msg := ""
//
//	for i := 0; i < len(domains); i++ {
//		if dr, ok := DataSync.Records[domains[i]]; ok {
//			msg += "DHash: " + hex.EncodeToString(domains[i][:])
//			msg += getDomain(dr)
//			msg += "\r\n"
//		}
//	}
//
//	return msg
//}
//
//func GetRecords(r string) string {
//	msg := ""
//	if r == "" {
//		for k, d := range DataSync.Records {
//			msg += "DHash: " + hex.EncodeToString(k[:])
//			msg += getDomain(d)
//			msg += "\r\n"
//		}
//
//		return msg
//	}
//
//	hash := Bas_Ethereum.GetHash(r)
//	msg += "DHash: " + r
//	if n, ok := DataSync.Records[hash]; !ok {
//		return "Domain Not Found"
//	} else {
//		msg = getDomain(n)
//	}
//
//	return msg
//
//}
//
//func getDomain(domain *DataSync.DomainRecord) string {
//	msg := ""
//
//	msg += fmt.Sprintf("   %-20s ", domain.GetName())
//	ip := domain.GetIPv4Addr()
//	msg += fmt.Sprintf("%-16s ", net.IPv4(ip[0], ip[1], ip[2], ip[3]).String())
//	msg += fmt.Sprintf("%-12s ", strconv.FormatInt(domain.GetExpire(), 10))
//	rare:="0"
//	if domain.GetIsRare() {
//		rare = "1"
//	}
//	msg += fmt.Sprintf("%-2s",rare)
//	open:="0"
//	if domain.GetOpenStatus(){
//		open = "1"
//	}
//	msg += fmt.Sprintf("%-2s",open)
//
//	return msg
//}
//
//func getDealString(deal *Market.Deal) string {
//	msg:=""
//
//	d:=api.GetRecord(deal.GetHash())
//	msg += fmt.Sprintf("%-20s",string(d.Name))
//	old := deal.GetFromOwner()
//	oldlen := len(old)
//	oldOwner := old[:4] + old[oldlen-4:]
//
//	own:=deal.GetOwner()
//	ownlen := len(own)
//
//	owner := own[:4] + own[ownlen-4:]
//
//	msg += fmt.Sprintf("%-9s",oldOwner)
//	msg += fmt.Sprintf("%-9s",owner)
//
//	msg += fmt.Sprintf("%-26s",deal.GetAGreedPrice().String())
//	//t,_:=DataSync.GetTimestamp(deal.BlockNumber)
//	msg += fmt.Sprintf("%-12s",strconv.FormatInt(Market.BlockNumnber2TimeStamp(deal.BlockNumber),10))
//
//	return msg
//
//}
//
//func GetDeal(domain string) string  {
//	msg := ""
//
//	domainhash:=Bas_Ethereum.Hash{}
//	if domain != ""{
//		domainhash = Bas_Ethereum.GetHash(domain)
//	}
//
//	for i:=0;i<len(Market.Sold);i++{
//		d:=&Market.Sold[i]
//		if domain == "" || (domain != "" && d.GetHash() == domainhash){
//			if msg != ""{
//				msg += "\r\n"
//			}
//			msg += getDealString(d)
//		}
//	}
//
//	return msg
//}
//
//func GetOrder(wallet string) string  {
//
//	msg := ""
//
//	if wallet != ""{
//		addr:=common.HexToAddress(wallet)
//		if m,ok:=Market.SellOrders[addr];!ok{
//			return "Not found"
//		}else{
//			msg = getOrderString(m)
//		}
//	}else{
//		for k,m:=range Market.SellOrders{
//			if msg != ""{
//				msg += "\r\n"
//			}
//			msg += k.String() + " "
//			if msg != ""{
//				msg += "\r\n"
//			}
//			msg += getOrderString(m)
//		}
//	}
//
//	return msg
//}
//
//func getOrderString(m map[Bas_Ethereum.Hash]*Market.SellOrder) string {
//	msg:=""
//
//	for k,v:=range m{
//		if msg != ""{
//			msg += "\r\n"
//		}
//		msg += k.String()+" "
//		d:=api.GetRecord(k)
//		if d == nil{
//			continue
//		}
//
//		msg += fmt.Sprintf("%-20s",string(d.Name))
//		msg += fmt.Sprintf("%-12s",strconv.FormatInt(Market.BlockNumnber2TimeStamp(v.BlockNumber),10))
//		msg += fmt.Sprintf("%-16s",v.GetPriceStr())
//	}
//
//	return msg
//
//}
