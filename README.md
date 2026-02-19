# 🛡️ CyFi: Forensic First Aid

### Securing the "Golden Hour" with Automated Live Response & Blockchain Integrity. CyFi is a next-generation digital forensics orchestration platform. It solves the "Golden Hour" crisis by automating the capture of volatile (RAM) and stable (Disk) evidence via a zero-dependency Go Agent and anchoring the evidence's cryptographic fingerprints onto a Blockchain Ledger to ensure an unbreakable chain of custody.

---

## Key Features  
**⚡ One-Click Acquisition:** A standalone, headless Go-based agent that performs RAM dumps and disk imaging silently.  
**📂 Forensic Core:** Orchestrates industry-standard tools like WinPmem (RAM) and FTK Imager CLI (Disk .E01).  
**⛓️ Immutable Chain of Custody:** Automatically anchors SHA-256 hashes of captured evidence to a private Ethereum-compatible blockchain (Hardhat/Polygon).  
**🔍 Integrity Verification:** A centralized dashboard to ingest USB forensic logs and verify evidence authenticity against the ledger.  
**📉 Zero-Footprint Design:** Minimized memory pollution during capture to preserve the integrity of the suspect machine.  

## 🛠️ Tech Stack

| Component | Technology |
| :--- | :--- |
| **Forensic Agent** | Golang (Compiled for Windows) |
| **Dashboard** | Streamlit (Python) |
| **Blockchain** | Solidity, Hardhat, Web3.py |
| **Data Integrity** | SHA-256 Hashing |
| **Binaries** | WinPmem, FTK Imager CLI |

---

## 📂 Project Structure
```text

├── agent/
│   ├── cyfi_agent.go       # Source code for the USB agent
│   └── tools/              # Place winpmem.exe and ftkimager.exe here
├── blockchain/
│   ├── contracts/          # Solidity Smart Contracts
│   └── artifacts/          # Compiled contract JSONs
├── dashboard/
│   └── app.py              # Streamlit Forensic Command Center
├── README.md
└── requirements.txt        # Python dependencies
```

     
## ⚙️ Setup & Installation  
1. The Blockchain (Hardhat)
   Ensure you have Node.js installed.
   ```text
   npm install --save-dev hardhat
   npx hardhat node
   ```
3. The Dashboard (Python)
   ```text
   pip install streamlit web3 psutil pandas
   streamlit run dashboard/app.py
   ```
4. The USB Agent (Go)
   Compile the Go agent to run without a terminal window
   ```text
   go build -ldflags "-H=windowsgui" agent/cyfi_agent.go
   ```

## 🖥️ Workflow  
**1. Capture:** Plug the CyFi USB into the suspect PC and run cyfi_agent.exe as Administrator.  
**2. Wait:** Upon completion, a native Windows alert will confirm the evidence is saved.  
**3. Ingest:** Plug the USB into the Forensic Workstation. The Streamlit dashboard will auto-detect the audit_log.txt.  
**4. Anchor:** Click "Initialize Vault" and then "Anchor to Ledger" to seal the hash on the blockchain.  
**5. Verify:** Use the Verification Engine to prove the evidence hasn't been tampered with since the second of capture.

## 🔮 Future Scope  
**IPFS Integration:** Decentralized storage for actual forensic images.  
**AI Triage:** Automatic scanning of RAM dumps for malware patterns.  
**Network Triage:** Remote heuristic monitoring to trigger auto-capture on suspicious nodes.  
  
**👥 Team CyFi**  
**Institution: Indira Gandhi Delhi Technical University for Women (IGDTUW)**  
**Track: Digital Forensics (WIEGNITE 3.0)**
