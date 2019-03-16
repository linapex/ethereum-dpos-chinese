
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 12:09:46</date>
//</624342666540224512>


package core

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"sync"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/log"
	"golang.org/x/crypto/ssh/terminal"
)

type CommandlineUI struct {
	in *bufio.Reader
	mu sync.Mutex
}

func NewCommandlineUI() *CommandlineUI {
	return &CommandlineUI{in: bufio.NewReader(os.Stdin)}
}

//readstring从stdin读取一行，从空格中剪裁if，强制
//非空性。
func (ui *CommandlineUI) readString() string {
	for {
		fmt.Printf("> ")
		text, err := ui.in.ReadString('\n')
		if err != nil {
			log.Crit("Failed to read user input", "err", err)
		}
		if text = strings.TrimSpace(text); text != "" {
			return text
		}
	}
}

//readpassword从stdin读取一行，从尾随的new
//行并返回它。输入将不会被回送。
func (ui *CommandlineUI) readPassword() string {
	fmt.Printf("Enter password to approve:\n")
	fmt.Printf("> ")

	text, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Crit("Failed to read password", "err", err)
	}
	fmt.Println()
	fmt.Println("-----------------------")
	return string(text)
}

//readpassword从stdin读取一行，从尾随的new
//行并返回它。输入将不会被回送。
func (ui *CommandlineUI) readPasswordText(inputstring string) string {
	fmt.Printf("Enter %s:\n", inputstring)
	fmt.Printf("> ")
	text, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Crit("Failed to read password", "err", err)
	}
	fmt.Println("-----------------------")
	return string(text)
}

//如果用户输入“是”，则confirm返回true，否则返回false
func (ui *CommandlineUI) confirm() bool {
	fmt.Printf("Approve? [y/N]:\n")
	if ui.readString() == "y" {
		return true
	}
	fmt.Println("-----------------------")
	return false
}

func showMetadata(metadata Metadata) {
	fmt.Printf("Request context:\n\t%v -> %v -> %v\n", metadata.Remote, metadata.Scheme, metadata.Local)
}

//approvetx提示用户确认请求签署交易
func (ui *CommandlineUI) ApproveTx(request *SignTxRequest) (SignTxResponse, error) {
	ui.mu.Lock()
	defer ui.mu.Unlock()
	weival := request.Transaction.Value.ToInt()
	fmt.Printf("--------- Transaction request-------------\n")
	if to := request.Transaction.To; to != nil {
		fmt.Printf("to:    %v\n", to.Original())
		if !to.ValidChecksum() {
			fmt.Printf("\nWARNING: Invalid checksum on to-address!\n\n")
		}
	} else {
		fmt.Printf("to:    <contact creation>\n")
	}
	fmt.Printf("from:  %v\n", request.Transaction.From.String())
	fmt.Printf("value: %v wei\n", weival)
	if request.Transaction.Data != nil {
		d := *request.Transaction.Data
		if len(d) > 0 {
			fmt.Printf("data:  %v\n", common.Bytes2Hex(d))
		}
	}
	if request.Callinfo != nil {
		fmt.Printf("\nTransaction validation:\n")
		for _, m := range request.Callinfo {
			fmt.Printf("  * %s : %s", m.Typ, m.Message)
		}
		fmt.Println()

	}
	fmt.Printf("\n")
	showMetadata(request.Meta)
	fmt.Printf("-------------------------------------------\n")
	if !ui.confirm() {
		return SignTxResponse{request.Transaction, false, ""}, nil
	}
	return SignTxResponse{request.Transaction, true, ui.readPassword()}, nil
}

//ApproveSignData提示用户确认请求签署数据
func (ui *CommandlineUI) ApproveSignData(request *SignDataRequest) (SignDataResponse, error) {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	fmt.Printf("-------- Sign data request--------------\n")
	fmt.Printf("Account:  %s\n", request.Address.String())
	fmt.Printf("message:  \n%q\n", request.Message)
	fmt.Printf("raw data: \n%v\n", request.Rawdata)
	fmt.Printf("message hash:  %v\n", request.Hash)
	fmt.Printf("-------------------------------------------\n")
	showMetadata(request.Meta)
	if !ui.confirm() {
		return SignDataResponse{false, ""}, nil
	}
	return SignDataResponse{true, ui.readPassword()}, nil
}

//approveexport提示用户确认导出加密帐户json
func (ui *CommandlineUI) ApproveExport(request *ExportRequest) (ExportResponse, error) {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	fmt.Printf("-------- Export Account request--------------\n")
	fmt.Printf("A request has been made to export the (encrypted) keyfile\n")
	fmt.Printf("Approving this operation means that the caller obtains the (encrypted) contents\n")
	fmt.Printf("\n")
	fmt.Printf("Account:  %x\n", request.Address)
//fmt.printf（“keyfile:\n%v\n”，request.file）
	fmt.Printf("-------------------------------------------\n")
	showMetadata(request.Meta)
	return ExportResponse{ui.confirm()}, nil
}

//approveImport提示用户确认导入账号json
func (ui *CommandlineUI) ApproveImport(request *ImportRequest) (ImportResponse, error) {
	ui.mu.Lock()
	defer ui.mu.Unlock()

	fmt.Printf("-------- Import Account request--------------\n")
	fmt.Printf("A request has been made to import an encrypted keyfile\n")
	fmt.Printf("-------------------------------------------\n")
	showMetadata(request.Meta)
	if !ui.confirm() {
		return ImportResponse{false, "", ""}, nil
	}
	return ImportResponse{true, ui.readPasswordText("Old password"), ui.readPasswordText("New password")}, nil
}

//批准提示用户确认列出帐户
//用户界面可以修改要列出的科目列表
func (ui *CommandlineUI) ApproveListing(request *ListRequest) (ListResponse, error) {

	ui.mu.Lock()
	defer ui.mu.Unlock()

	fmt.Printf("-------- List Account request--------------\n")
	fmt.Printf("A request has been made to list all accounts. \n")
	fmt.Printf("You can select which accounts the caller can see\n")
	for _, account := range request.Accounts {
		fmt.Printf("\t[x] %v\n", account.Address.Hex())
	}
	fmt.Printf("-------------------------------------------\n")
	showMetadata(request.Meta)
	if !ui.confirm() {
		return ListResponse{nil}, nil
	}
	return ListResponse{request.Accounts}, nil
}

//ApproveWaccount提示用户确认创建新帐户，并显示给调用方
func (ui *CommandlineUI) ApproveNewAccount(request *NewAccountRequest) (NewAccountResponse, error) {

	ui.mu.Lock()
	defer ui.mu.Unlock()

	fmt.Printf("-------- New Account request--------------\n")
	fmt.Printf("A request has been made to create a new. \n")
	fmt.Printf("Approving this operation means that a new Account is created,\n")
	fmt.Printf("and the address show to the caller\n")
	showMetadata(request.Meta)
	if !ui.confirm() {
		return NewAccountResponse{false, ""}, nil
	}
	return NewAccountResponse{true, ui.readPassword()}, nil
}

//ShowError向用户显示错误消息
func (ui *CommandlineUI) ShowError(message string) {

	fmt.Printf("ERROR: %v\n", message)
}

//ShowInfo向用户显示信息消息
func (ui *CommandlineUI) ShowInfo(message string) {
	fmt.Printf("Info: %v\n", message)
}

func (ui *CommandlineUI) OnApprovedTx(tx ethapi.SignTransactionResult) {
	fmt.Printf("Transaction signed:\n ")
	spew.Dump(tx.Tx)
}

func (ui *CommandlineUI) OnSignerStartup(info StartupInfo) {

	fmt.Printf("------- Signer info -------\n")
	for k, v := range info.Info {
		fmt.Printf("* %v : %v\n", k, v)
	}
}

