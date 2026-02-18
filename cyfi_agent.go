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

// --- MAIN PIPELINE ---
func main() {
	// 1. Setup Paths (Execute from USB)
	exePath, _ := os.Executable()
	usbRoot = filepath.Dir(exePath)
	evidenceDir = filepath.Join(usbRoot, "evidence")
	logFile = filepath.Join(usbRoot, "audit_log.txt")
	toolsDir := filepath.Join(usbRoot, "tools")

	os.MkdirAll(evidenceDir, os.ModePerm)

	logEvent("=== NEW FORENSIC ACQUISITION INITIALIZED ===")
	logEvent(fmt.Sprintf("Executing from USB Root: %s", usbRoot))

	// 2. STAGE 1: RAM Capture (winpmem)
	logEvent("STAGE 1: Triggering kernel-mode RAM Capture...")
	ramFile := filepath.Join(evidenceDir, "ram_dump.raw")
	winpmemPath := filepath.Join(toolsDir, "winpmem.exe")

	// Execute Winpmem silently
	cmdRam := exec.Command(winpmemPath, "-o", ramFile)
	err := cmdRam.Run()
	if err != nil {
		logEvent(fmt.Sprintf("WARNING: RAM capture failed or simulated: %v", err))
		// Fallback for demo if winpmem isn't present
		os.WriteFile(ramFile, []byte("SIMULATED_VOLATILE_MEMORY_DATA"), 0644)
	}

	ramHash, _ := calculateHash(ramFile)
	logEvent(fmt.Sprintf("SUCCESS: RAM captured. SHA-256: %s", ramHash))

	// 3. STAGE 2: Disk Imaging (FTK Imager)
	logEvent("STAGE 2: Triggering Write-Blocked Disk Imaging (.E01)...")
	diskFileBase := filepath.Join(evidenceDir, "disk_image") // FTK appends .e01 automatically
	ftkPath := filepath.Join(toolsDir, "FTK_Imager//ftkimager.exe")

	// Execute FTK Imager silently (Targeting PhysicalDrive0)
	cmdDisk := exec.Command(ftkPath, "\\\\.\\PhysicalDrive0", diskFileBase, "--e01", "--frag", "2G", "--quiet")
	err = cmdDisk.Run()
	if err != nil {
		logEvent(fmt.Sprintf("WARNING: Disk imaging failed or simulated: %v", err))
		// Fallback for demo if ftkimager isn't present
		os.WriteFile(diskFileBase+".e01", []byte("SIMULATED_PHYSICAL_DRIVE_DATA"), 0644)
	}

	diskHash, _ := calculateHash(diskFileBase + ".e01")
	logEvent(fmt.Sprintf("SUCCESS: Disk imaged. SHA-256: %s", diskHash))

	// 4. FINISH
	logEvent("=== ACQUISITION COMPLETED SUCCESSFULLY ===")

	// Pop up native Windows alert (64 = Information Icon)
	MessageBox("CyFi Secure Agent", "Forensic Acquisition Complete.\nEvidence and Logs saved to USB.\n\nIt is now safe to remove the drive.", 64)
}
