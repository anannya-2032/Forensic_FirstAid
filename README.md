# CyFi- Forensic First Aid: Live-to-Disk Pipeline

## Overview

*Forensic First Aid* (by Team *CyFi) is a unified, "one-click" USB workflow designed to solve the "Golden Hour" crisis in digital forensics. When a security breach occurs, critical evidence in RAM is often lost due to reboots or anti-forensic wipers. This project automates the acquisition of volatile and non-volatile data, securing every artifact with a **SHA-256 digital fingerprint* and an *Immutable Audit Log* anchored to a blockchain.

## The Problem: The "Golden Hour" Crisis

* *The Volatility Gap:* RAM, encryption keys, and live malware traces are lost if a system is powered down before capture.


* *Speed-Soundness Trade-off:* Manual forensic setups are slow, giving "anti-forensic" wipers time to destroy evidence.


* *Chain of Custody Fragility:* Standard logs are easily manipulated, leading to trust deficits in judicial reviews.



## Key Features

* *Near-Zero Footprint:* Built as a *Go-based static binary* to execute without heavy interpreters or dependencies on the suspect host.


* *Blockchain-Anchored Integrity:* Uses *Solidity* smart contracts to create a decentralized, immutable record of evidence timestamps.


* *Automated 6-Stage Pipeline:* Denies malware the "window of opportunity" by rapidly bridging the gap between RAM and disk acquisition.



---

## 6-Stage Automated Pipeline

1. *RAM State:* Freezes volatile processes and active connections into a "live snapshot".


2. *Registry Hive:* Extracts the system’s "Binary DNA" (HKLM) to document hardware, users, and persistent settings.


3. *Network State:* Logs active IP connections to identify unauthorized remote access or data exfiltration.


4. *DNS Cache:* Recovers visited domain history, bypassing "cleared" browser histories.


5. *User Intent:* Captures PowerShell command history—the "Smoking Gun" for manual hacking.


6. *Security Logs:* Exports the Windows Event Log to create a chronological timeline of logins and privilege changes.



---

## System Architecture

The system consists of two primary environments:

1. *Forensic Host (Suspect PC):* Runs the cyfi_agent.exe (Go Agent) from a USB drive to write evidence.raw and audit_log.txt.


2. *Forensic Workstation:* A Streamlit-based dashboard used for log ingestion, SHA-256 verification, and blockchain anchoring.



---

## Tech Stack

| Layer | Technology | Purpose |
| --- | --- | --- |
| *User Interface* | Streamlit (Python) | No-code investigator dashboard.|
| *Logic/Brain* | Go (Golang) | Lightweight, portable execution with zero dependencies.|
| *Integrity Layer* | Solidity | Immutable, decentralized Chain of Custody record.|
| *Environment* | Hardhat / Ganache | Local blockchain node for anchoring evidence.|
| *Communication* | JSON | Linking local tools to the blockchain vault.|

---

## Future Scope

* *Multi-Cloud Storage (IPFS):* Moving from local file storage to a fully decentralized InterPlanetary File System.


* *AI-Powered Analysis:* An AI/ML layer to automatically scan massive images for malware signatures or hidden partitions.


* *Insurance Integration:* Automating "insurance-ready" reports that meet specific evidentiary criteria for claim payouts.



## Team CyFi

Developed at *Indira Gandhi Delhi Technical University for Women (IGDTUW)* for *WIEGNITE 3.0*.
