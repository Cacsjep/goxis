### Step 1: Access the Switch
Connect to your Cisco switch via console cable or SSH/telnet, and enter the privileged EXEC mode:

enable

### Step 2: Enter Configuration Mode
configure terminal

### Step 3: Create VLANs
Create a VLAN, for example, VLAN 10:

vlan 10
name VLAN10
exit

Repeat the above commands to create additional VLANs if necessary.

### Step 4: Assign Ports to VLANs
Configure ports 1 to 16 for VLAN 10. Ports will be set to access mode, and they will send untagged traffic:

interface range fastEthernet 0/1 - 0/16
switchport mode access
switchport access vlan 10
exit

### Step 5: Configure the Uplink Port
Set the uplink port (e.g., port 24) to trunk mode. This port will carry all VLANs. By default, a trunk port carries all VLANs, but you can specify VLANs if needed:

interface fastEthernet 0/24
switchport mode trunk
switchport trunk allowed vlan all
exit

In this setup, the uplink port will tag VLANs, but you can configure it to handle native VLAN traffic as untagged. For instance, if you want VLAN 10 to be untagged:

switchport trunk native vlan 10

mathematica
Copy code

### Step 6: Save the Configuration
Save your configuration to ensure it persists across reboots:

write memory

**Notes:**
- Replace `fastEthernet` with `gigabitEthernet` or another interface type if your switch uses different interface naming.
- Ensure your switch’s IOS version supports the commands used. Commands can vary slightly between models and IOS versions.
- Configure additional settings such as VTP, STP as required for your network environment.

This set of instructions will set up VLAN 10 across ports 1 to 16 as untagged, and the uplink port