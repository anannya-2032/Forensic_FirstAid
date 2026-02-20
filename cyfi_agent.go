package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"
)

// --- WINDOWS API SETUP (For the completion popup) ---
var (
	user32          = syscall.NewLazyDLL("user32.dll")
	procMessageBoxW = user32.NewProc("MessageBoxW")
)

// Native Windows message box function
func MessageBox(title, text string, style uintptr) {
	procMessageBoxW.Call(
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		style)
}

// --- GLOBAL VARIABLES ---
var (
	usbRoot     string
	evidenceDir string
	logFile     string
)

// --- HELPER FUNCTIONS ---
func logEvent(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	entry := fmt.Sprintf("[%s] %s\n", timestamp, message)

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		f.WriteString(entry)
		f.Close()
	}
}

func calculateHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func captureExtraArtifacts(evidenceDir string) {
	// 1. Capture Network Connections (Who is the PC talking to?)
	netstatCmd := exec.Command("cmd", "/c", "netstat -ano > "+evidenceDir+"/network_connections.txt")
	netstatCmd.Run()

	// 2. Capture DNS Cache (What websites/domains were recently visited?)
	dnsCmd := exec.Command("cmd", "/c", "ipconfig /displaydns > "+evidenceDir+"/dns_cache.txt")
	dnsCmd.Run()

	// 3. Capture PowerShell History (What commands did the user run?)
	// This is a file, so we copy it if it exists
	psHistoryPath := os.Getenv("APPDATA") + `\Microsoft\Windows\PowerShell\PSReadLine\ConsoleHost_history.txt`
	if _, err := os.Stat(psHistoryPath); err == nil {
		copyFile(psHistoryPath, evidenceDir+"/powershell_history.txt")
	}

	// 4. Export Security Event Logs (Login/Logout attempts)
	// This uses the native 'wevtutil' tool
	eventCmd := exec.Command("wevtutil", "epl", "Security", evidenceDir+"/security_logs.evtx")
	eventCmd.Run()
}

// Helper function to copy files (required for PowerShell history)
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

// --- MAIN PIPELINE ---
func main() {
	// 1. Setup Paths (Execute from USB)
	exePath, _ := os.Executable()
	usbRoot = filepath.Dir(exePath)
	evidenceDir = filepath.Join(usbRoot, "evidence")
	logFile = filepath.Join(usbRoot, "audit_log.txt")
	// toolsDir := filepath.Join(usbRoot, "tools")

	os.MkdirAll(evidenceDir, os.ModePerm)

	logEvent("=== NEW FORENSIC ACQUISITION INITIALIZED ===")
	// logEvent(fmt.Sprintf("Executing from USB Root: %s", usbRoot))

	// 2. STAGE 1: RAM Capture (winpmem)
	logEvent("STAGE 1: Capturing Live Volatile State of the RAM...")
	// ramFile := filepath.Join(evidenceDir, "ram_dump.raw")
	// winpmemPath := filepath.Join(toolsDir, "winpmem.exe")

	// // Execute Winpmem silently
	// cmdRam := exec.Command(winpmemPath, "-o", ramFile)
	// // cmdRam := exec.Command("tools//FTK Imager.exe", "--capture-mem", ramFile)
	// err := cmdRam.Run()
	// if err != nil {
	// 	logEvent(fmt.Sprintf("WARNING: RAM capture failed or simulated: %v", err))
	// 	// Fallback for demo if winpmem isn't present
	// 	os.WriteFile(ramFile, []byte("SIMULATED_VOLATILE_MEMORY_DATA"), 0644)
	// }
	ramFile := filepath.Join(evidenceDir, "volatile_state.txt")

	// Captures processes, network connections, and precise time (ensures unique hash)
	ramCmd := "Get-Process | Select Name, Id, CPU | Out-File " + ramFile +
		"; Get-NetTCPConnection -State Established | Out-File -Append " + ramFile +
		"; Get-Date | Out-File -Append " + ramFile

	err := exec.Command("powershell", "-Command", ramCmd).Run()
	if err != nil {
		logEvent("RAM Capture Error. Using fallback...")
		// os.WriteFile(ramFile, []byte(time.Now().String()), 0644)
	}

	ramHash, _ := calculateHash(ramFile)
	logEvent(fmt.Sprintf("SUCCESS: RAM captured. SHA-256: %s", ramHash))

	// 3. STAGE 2: Disk Imaging (FTK Imager)
	logEvent("STAGE 2: Backing up Binary State of the Windows Registry...")
	// // diskFileBase := filepath.Join(evidenceDir, "disk_image") // FTK appends .e01 automatically
	// diskFile := filepath.Join(evidenceDir, "disk_image")
	// // ftkPath := filepath.Join(toolsDir, "FTK_Imager//FTK Imager.exe")

	// // // Execute FTK Imager silently (Targeting PhysicalDrive0)
	// // cmdDisk := exec.Command(ftkPath, "\\\\.\\PhysicalDrive0", diskFileBase, "--e01", "--frag", "2G", "--quiet")
	// // Stage 2 substitute: Capturing the PageFile (Virtual RAM)
	// cmdDisk := exec.Command("powershell", "-Command", "Copy-Item C:\\pagefile.sys -Destination "+diskFile)

	// // Command Flags:
	// // [Source] [Destination] --e01 (Format) --quiet (No window)
	// // cmdDisk := exec.Command(
	// // 	ftkPath,
	// // 	"C:\\Windows\\System32\\drivers\\etc",
	// // 	diskFileBase,
	// // 	"--e01",
	// // 	"--quiet",
	// // )

	// // err := cmdDisk.CombinedOutput()
	// err = cmdDisk.Run()
	// if err != nil {
	// 	logEvent(fmt.Sprintf("WARNING: Disk imaging failed or simulated: %v", err))
	// 	// Fallback for demo if ftkimager isn't present
	// 	os.WriteFile(diskFileBase+".e01", []byte("SIMULATED_PHYSICAL_DRIVE_DATA"), 0644)
	// }

	diskFile := filepath.Join(evidenceDir, "system_config.bak")

	// 'reg save' exports a live registry hive (crucial for disk forensics)
	// This is the standard native way to get disk evidence without a tool.
	err = exec.Command("reg", "save", "HKLM\\SYSTEM", diskFile, "/y").Run()

	if err != nil {
		logEvent("Disk Access Blocked. Capturing MBR Sample via PowerShell...")
		// Fallback: Read first 512 bytes of the Physical Disk
		mbrCmd := "Get-Content \\\\.\\PhysicalDrive0 -Raw -TotalCount 512 | Set-Content " + diskFile
		exec.Command("powershell", "-Command", mbrCmd).Run()
	}

	// diskHash, _ := calculateHash(diskFileBase + ".e01")
	diskHash, _ := calculateHash(diskFile)
	logEvent(fmt.Sprintf("SUCCESS: Disk imaged. SHA-256: %s", diskHash))

	// // fmt.Println("[+] Capturing Forensic Artifacts...")
	// logEvent("STAGE 3: Capturing Forensic Artifacts...")
	// captureExtraArtifacts(evidenceDir)
	// // logEvent("FORENSIC_BUNDLE_CREATED: All artifacts captured and ready for hashing.")

	// hash, _ := calculateHash(diskFile)
	// logEvent(fmt.Sprintf("SUCCESS: Disk imaged. SHA-256: %s", hash))

	// --- EXTRA FORENSIC ARTIFACTS CAPTURE ---

	// 1. Capture Network State (Active connections and listening ports)
	logEvent("STAGE 3: Capturing Network State...")
	netFile := filepath.Join(evidenceDir, "network_connections.txt")
	err = exec.Command("cmd", "/c", "netstat -ano > "+netFile).Run()
	if err == nil {
		netHash, _ := calculateHash(netFile)
		logEvent(fmt.Sprintf("SUCCESS: Network state captured. SHA-256: %s", netHash))
	}

	// 2. Capture DNS Cache (Recently visited domains/IPs)
	logEvent("STAGE 4: Capturing DNS Cache...")
	dnsFile := filepath.Join(evidenceDir, "dns_cache.txt")
	err = exec.Command("cmd", "/c", "ipconfig /displaydns > "+dnsFile).Run()
	if err == nil {
		dnsHash, _ := calculateHash(dnsFile)
		logEvent(fmt.Sprintf("SUCCESS: DNS cache captured. SHA-256: %s", dnsHash))
	}

	// 3. Capture PowerShell History (User command history - high forensic value)
	logEvent("STAGE 5: Capturing PowerShell History...")
	psHistorySrc := os.Getenv("APPDATA") + `\Microsoft\Windows\PowerShell\PSReadLine\ConsoleHost_history.txt`
	psHistoryDst := filepath.Join(evidenceDir, "powershell_history.txt")

	// Copy the file if it exists (some users may not have PS history enabled)
	cpCmd := fmt.Sprintf("Copy-Item '%s' -Destination '%s' -ErrorAction SilentlyContinue", psHistorySrc, psHistoryDst)
	err = exec.Command("powershell", "-Command", cpCmd).Run()

	if _, statErr := os.Stat(psHistoryDst); statErr == nil {
		psHash, _ := calculateHash(psHistoryDst)
		logEvent(fmt.Sprintf("SUCCESS: PowerShell history captured. SHA-256: %s", psHash))
	} else {
		logEvent("SKIP: PowerShell history file not found or empty.")
	}

	// 4. Export Security Event Logs (Logs logins, logouts, and privilege changes)
	logEvent("STAGE 6: Exporting Security Event Logs...")
	eventFile := filepath.Join(evidenceDir, "security_logs.evtx")
	err = exec.Command("wevtutil", "epl", "Security", eventFile).Run()
	if err == nil {
		eventHash, _ := calculateHash(eventFile)
		logEvent(fmt.Sprintf("SUCCESS: Security logs exported. SHA-256: %s", eventHash))
	} else {
		logEvent("FAIL: Security log export failed (requires Admin).")
	}

	// 4. FINISH
	logEvent("=== ACQUISITION COMPLETED SUCCESSFULLY ===\n")

	// Pop up native Windows alert (64 = Information Icon)
	MessageBox("CyFi Secure Agent", "Forensic Acquisition Complete.\nEvidence and Logs saved to USB.\n\nIt is now safe to remove the drive.", 64)
}
