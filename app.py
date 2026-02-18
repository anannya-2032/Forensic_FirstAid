import streamlit as st
from web3 import Web3
import json
import os
import time
import hashlib
import datetime
import pandas as pd
import psutil

# 1. PAGE SETUP
st.set_page_config(page_title="CyFi: Forensic Command Center", page_icon="🛡️", layout="wide")

# 2. INITIALIZE SESSION STATE
if 'vault_address' not in st.session_state:
    st.session_state.vault_address = "0x5FbDB2315678afecb367f032d93F642f64180aa3"
if 'history' not in st.session_state:
    st.session_state.history = []
if 'current_hash' not in st.session_state:
    st.session_state.current_hash = ""

# 3. BLOCKCHAIN CONNECTION & ABI
w3 = Web3(Web3.HTTPProvider("http://127.0.0.1:8545"))
ABI = [
    {"inputs":[{"name":"_fileHash","type":"string"}],"name":"anchorEvidence","outputs":[],"stateMutability":"nonpayable","type":"function"},
    {"inputs":[{"name":"_fileHash","type":"string"}],"name":"verifyEvidence","outputs":[{"name":"","type":"uint256"}],"stateMutability":"view","type":"function"}
]

# --- HELPER: USB DETECTION ---
def auto_find_usb():
    for part in psutil.disk_partitions():
        try:
            log_path = os.path.join(part.mountpoint, "audit_log.txt")
            if os.path.exists(log_path):
                return part.mountpoint
        except:
            continue
    return None

# --- SIDEBAR ---
st.sidebar.title("🛡️ CyFi System Status")
if w3.is_connected():
    st.sidebar.success("⛓️ Ledger: CONNECTED")
    st.sidebar.info(f"📍 Vault Address: \n{st.session_state.vault_address}")
else:
    st.sidebar.error("❌ Ledger: OFFLINE")
    st.sidebar.warning("Run 'npx hardhat node' in a terminal")
    if st.sidebar.button("🔄 Refresh Connection"):
        st.rerun()

# --- MAIN UI ---
st.title("🛡️ CyFi: Digital Evidence Vault")
st.markdown("---")

# ROW 1: USB INGESTION (HYBRID)
st.subheader("🔌 Evidence Source Selection")
col_auto, col_manual = st.columns([1, 1])

detected_path = auto_find_usb()

with col_auto:
    if detected_path:
        st.success(f"✅ Auto-Detected USB at: {detected_path}")
        final_path = detected_path
    else:
        st.warning("🔍 No USB with 'audit_log.txt' detected automatically.")
        if st.button("🔄 Rescan Drives"):
            st.rerun()

with col_manual:
    manual_path = st.text_input("OR Enter Drive/Folder Path manually:", placeholder="e.g. E:\\ or C:\\Forensics\\USB_Dump")
    if manual_path:
        final_path = manual_path
    elif detected_path:
        final_path = detected_path
    else:
        final_path = None

st.markdown("---")

# ROW 2: LOG & FILE DISPLAY
if final_path:
    log_file_path = os.path.join(final_path, "audit_log.txt")
    evidence_folder_path = os.path.join(final_path, "evidence")
    
    c1, c2 = st.columns(2)
    
    with c1:
        st.markdown("### 📜 USB Audit Log")
        if os.path.exists(log_file_path):
            with open(log_file_path, "r") as f:
                content = f.read()
            st.text_area("USB Log Data (Copy Hash from here)", value=content, height=200)
        else:
            st.error(f"audit_log.txt not found at {final_path}")

    with c2:
        st.markdown("### 📂 Evidence Files")
        if os.path.exists(evidence_folder_path):
            files = os.listdir(evidence_folder_path)
            if files:
                st.write(files)
            else:
                st.info("Evidence folder exists but is empty.")
        else:
            st.error("Evidence folder not found.")
else:
    st.info("Please connect the USB or enter the manual path to view logs and evidence.")

st.markdown("---")

# ROW 3: BLOCKCHAIN ANCHORING & ADMIN (From your previous code)
col1, col2 = st.columns(2)

with col1:
    st.subheader("⚙️ Vault Admin")
    if st.button("🚀 Initialize New Vault", use_container_width=True):
        try:
            path = "./artifacts/contracts/CyFiVault.sol/CyFiVault.json"
            with open(path, "r") as f: artifact = json.load(f)
            factory = w3.eth.contract(abi=ABI, bytecode=artifact['bytecode'])
            tx = factory.constructor().transact({'from': w3.eth.accounts[0]})
            receipt = w3.eth.wait_for_transaction_receipt(tx)
            st.session_state.vault_address = receipt.contractAddress
            st.success("Vault Deployed!")
            st.rerun()
        except Exception as e:
            st.error(f"Deployment Failed: {e}")

with col2:
    st.subheader("⛓️ Blockchain Anchoring")
    # This is where you paste the hash copied from the audit log above
    manual_hash = st.text_input("Paste Hash to Anchor", value=st.session_state.current_hash)
    if st.button("🔒 Anchor to Ledger", use_container_width=True):
        if not st.session_state.vault_address:
            st.error("Please Initialize a Vault first!")
        else:
            try:
                contract = w3.eth.contract(address=st.session_state.vault_address, abi=ABI)
                tx = contract.functions.anchorEvidence(manual_hash).transact({'from': w3.eth.accounts[0]})
                w3.eth.wait_for_transaction_receipt(tx)
                
                new_entry = {
                    "Time": datetime.datetime.now().strftime("%H:%M:%S"),
                    "Evidence Hash": f"{manual_hash[:15]}...",
                    "TX ID": f"{tx.hex()[:15]}...",
                    "Status": "✅ Sealed"
                }
                st.session_state.history.insert(0, new_entry)
                st.success("Immutably Sealed to Ledger!")
            except Exception as e:
                st.error(f"Anchoring Error: {e}")

# ROW 4: AUDIT LOG DISPLAY
st.markdown("---")
st.subheader("📜 Blockchain Evidence History")
if st.session_state.history:
    st.table(pd.DataFrame(st.session_state.history))
else:
    st.info("The blockchain ledger log is currently empty.")

# ROW 5: VERIFICATION ENGINE
st.markdown("---")
st.subheader("🕵️ Verification Engine")
check_hash = st.text_input("Paste Hash to Verify Integrity")

if st.button("🔍 Check Ledger Status"):
    if not w3.is_connected():
        st.error("Blockchain is not connected. Run 'npx hardhat node'.")
    elif not st.session_state.vault_address:
        st.error("No vault address found. Click 'Initialize New Vault' first.")
    else:
        try:
            # Check if contract exists on-chain
            code = w3.eth.get_code(st.session_state.vault_address)
            if code == b'' or code.hex() == "0x":
                st.error("Vault not found. Please 'Initialize New Vault' to redeploy.")
            else:
                verify_contract = w3.eth.contract(address=st.session_state.vault_address, abi=ABI)
                timestamp = verify_contract.functions.verifyEvidence(check_hash).call()
                
                if timestamp > 0:
                    dt = datetime.datetime.fromtimestamp(timestamp).strftime('%Y-%m-%d %H:%M:%S')
                    st.success(f"✅ VERIFIED: This evidence was anchored on {dt}")
                    st.balloons()
                else:
                    st.error("❌ UNVERIFIED: This hash does not exist in the current ledger.")
        except Exception as e:
            st.error(f"Verification Error: {e}")