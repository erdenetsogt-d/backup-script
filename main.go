package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	mikrotikIP := os.Getenv("MIKROTIK_IP")
	username := os.Getenv("MIKROTIK_USERNAME")
	password := os.Getenv("MIKROTIK_PASSWORD")
	sshPortstr := os.Getenv("PORT")
	c_ip := os.Getenv("FORTI_IP")
	c_uname := os.Getenv("FORTI_USERNAME")
	c_pass := os.Getenv("FORTI_PASSWORD")
	c_portstr := os.Getenv("FORTI_PORT")
	sshPort, _ := strconv.Atoi(sshPortstr)
	c_port, _ := strconv.Atoi(c_portstr)
	ftp_ip := os.Getenv("FTP_HOST")
	ftp_user := os.Getenv("FTP_USERNAME")
	ftp_pass := os.Getenv("FTP_PASS")
	ftp_path := os.Getenv("FTP_PATH")
	args := os.Args

	if len(os.Args) < 2 {
		fmt.Println("Please provide a parameter: mikrotik or fortinet. for example  ./backupper mikrotik or ./backupper fortinet")
		os.Exit(1)
	}
	param := args[1]
	fmt.Println("Parameter: " + param)
	if param == "mikrotik" {

		// Use the current timestamp to create a unique backup file name
		timestamp := time.Now().Format("20060102150405")
		backupFileName := fmt.Sprintf("mikrotik_config_%s.rsc", timestamp)
		backupFilePath := fmt.Sprintf("backup/%s", backupFileName)

		// Run the backup command using SSH
		cmd := exec.Command("sshpass", "-p", password, "ssh", "-p", fmt.Sprint(sshPort), "-o", "StrictHostKeyChecking=no", username+"@"+mikrotikIP, fmt.Sprintf(`/export file="%s"`, backupFileName))

		// Capture stdout and stderr
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatalf("Error creating StdoutPipe: %v", err)
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			log.Fatalf("Error creating StderrPipe: %v", err)
		}

		err = cmd.Start()
		if err != nil {
			log.Fatalf("Error starting command: %v", err)
		}

		// Read stdout and stderr
		stdoutBytes, err := io.ReadAll(stdout)
		if err != nil {
			log.Fatalf("Error reading stdout: %v", err)
		}
		stderrBytes, err := io.ReadAll(stderr)
		if err != nil {
			log.Fatalf("Error reading stderr: %v", err)
		}

		err = cmd.Wait()
		if err != nil {
			log.Printf("Command finished with error: %v", err)
			log.Printf("Stdout: %s", stdoutBytes)
			log.Printf("Stderr: %s", stderrBytes)
			log.Fatal("Failed to run MikroTik export command")
		} else {
			log.Println("MikroTik export command completed successfully")
			scpCmd := exec.Command("sshpass", "-p", password, "scp", "-P", fmt.Sprint(sshPort), "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", username+"@"+mikrotikIP+":"+backupFileName, backupFilePath)

			scpOutput, err := scpCmd.CombinedOutput()
			if err != nil {
				log.Printf("Error downloading file with scp: %v", err)
				log.Printf("Scp output: %s", scpOutput)
				log.Fatal("Failed to download the file")
			} else {
				log.Println("File downloaded successfully")
				rmCmd := exec.Command("sshpass", "-p", password, "ssh", "-p", fmt.Sprint(sshPort), "-o", "StrictHostKeyChecking=no", username+"@"+mikrotikIP, "rm", backupFileName)
				rmOutput, err := rmCmd.CombinedOutput()
				if err != nil {
					log.Printf("Error removing file on MikroTik device: %v", err)
					log.Printf("Rm output: %s", rmOutput)
					log.Fatal("Failed to remove the file on MikroTik device")
				} else {
					log.Println("File removed successfully on MikroTik device")
				}
			}
		}
	}
	if param == "fortinet" {
		timestamp := time.Now().Format("20060102150405")
		backupFileName := fmt.Sprintf("fortinet_config_%s.cfg", timestamp)
		ftp_path = ftp_path + backupFileName

		// Replace the following command with the actual Fortinet backup command
		fortinetBackupCommand := ("execute backup full-config sftp " + "'" + ftp_path + "'" + " " + ftp_ip + " " + ftp_user + " " + ftp_pass)
		fmt.Println(fortinetBackupCommand)
		cmd := exec.Command("sshpass", "-p", c_pass, "ssh", "-p", fmt.Sprint(c_port), "-o", "StrictHostKeyChecking=no", fmt.Sprintf("%s@%s", c_uname, c_ip), fortinetBackupCommand)

		err := cmd.Run()
		if err != nil {
			fmt.Printf("Error running Fortinet backup command: %v\n", err)
			return
		}

		fmt.Printf("Fortinet backup completed successfully. Backup file: %s\n", ftp_path)
	}
}
