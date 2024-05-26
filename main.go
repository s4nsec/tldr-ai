package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// GPTRequest represents the request payload for the GPT API.
type GPTRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

// Message represents a message in the conversation.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// GPTResponse represents the response from the GPT API.
type GPTResponse struct {
	Choices []Choice `json:"choices"`
}

// Choice represents a single choice in the GPT API response.
type Choice struct {
	Message Message `json:"message"`
}

func askGPT(utilName string) {
	client := resty.New()
	messageContent := `
    tldr is a utility that is to give brief information about a provided cli utility.
    For example, if execute "tldr nmap" in the terminal, I get:

    Successfully updated local database

    nmap

    Network exploration tool and security/port scanner.
    Some features (e.g. SYN scan) activate only when nmap is run with root privileges.
    More information: <https://nmap.org/book/man.html>.

    - Scan the top 1000 ports of a remote host with various [v]erbosity levels:
        nmap -v1|2|3 ip_or_hostname

    - Run a ping sweep over an entire subnet or individual hosts very aggressively:
        nmap -T5 -sn 192.168.0.0/24|ip_or_hostname1,ip_or_hostname2,...

    - Enable OS detection, version detection, script scanning, and traceroute:
        sudo nmap -A ip_or_hostname1,ip_or_hostname2,...

    - Scan a specific list of ports (use -p- for all ports from 1 to 65535):
        nmap -p port1,port2,... ip_or_host1,ip_or_host2,...

    - Perform service and version detection of the top 1000 ports using default NSE scripts, writing results (-oA) to output files:
        nmap -sC -sV -oA top-1000-ports ip_or_host1,ip_or_host2,...

    - Scan target(s) carefully using default and safe NSE scripts:
        nmap --script "default and safe" ip_or_host1,ip_or_host2,...

    - Scan for web servers running on standard ports 80 and 443 using all available http-* NSE scripts:
        nmap --script "http-*" ip_or_host1,ip_or_host2,... -p 80,443

    - Attempt evading IDS/IPS detection by using an extremely slow scan (-T0), decoy source addresses (-D), [f]ragmented packets, random data and other methods:
        sudo nmap -T0 -D decoy_ip1,decoy_ip2,... --source-port 53 -f --data-length 16 -Pn ip_or_host


    Provide me similar output for the utility: ` + utilName + ` and don't say anything else.`

	request := GPTRequest{
		Model: "gpt-4-turbo",
		Messages: []Message{
			{
				Role:    "user",
				Content: messageContent,
			},
		},
		Temperature: 0.7,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to marshal request body")
	}

    openaiAPIKey := os.Getenv("OPENAI_API_KEY")
    if openaiAPIKey == "" {
        log.Fatal().Msg("OPENAI_API_KEY environment variable not set")
    }

	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+ openaiAPIKey).
		SetBody(requestBody).
		Post("https://api.openai.com/v1/chat/completions")

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to make request to GPT API")
	}

	var gptResponse GPTResponse
	if err := json.Unmarshal(response.Body(), &gptResponse); err != nil {
		log.Fatal().Err(err).Msg("Failed to unmarshal response body")
	}

	if len(gptResponse.Choices) > 0 {
		message := gptResponse.Choices[0].Message
		fmt.Print(message.Content)
	} else {
		log.Warn().Msg("No response from GPT API")
	}
}

func tldrExists() {
	_, err := exec.LookPath("tldr")
	if err != nil {
		operatingSystem := runtime.GOOS
		switch operatingSystem {
		case "windows":
			exec.Command("choco", "install", "tldr").Run()
		case "linux":
			exec.Command("sudo", "apt", "install", "tldr").Run()
		case "darwin":
			exec.Command("brew", "install", "tldr").Run()
		default:
			log.Fatal().Msg("Operating system not supported")
		}
	}
}

func main() {
	// set log level
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// set log format
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// abort if no cli arguments
	if len(os.Args) < 2 {
		log.Fatal().Msgf("Usage: %s <command> <param>\n", os.Args[0])
	}

	// check if the command exists

	// get cli arguments
	cmdToRun := "tldr"
	paramToCmd := os.Args[1]

	// run cli Command
	cmd := exec.Command(cmdToRun, paramToCmd)
	output, err:= cmd.Output()

    if err != nil && !strings.Contains(err.Error(), "executable file not found") {
        tldrExists()
    }

	if !strings.Contains(string(output), "This page doesn't exist yet!") {
		fmt.Println(string(output))
	}

	if err != nil {
		askGPT(paramToCmd)
	}

}
