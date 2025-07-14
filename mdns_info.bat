@echo off
echo.
echo  About mDNS Multicast Interface Warnings on Windows
echo ===================================================
echo.
echo You may see warnings like:
echo [WARN] mdns: Failed to set multicast interface: no such interface
echo.
echo  What this means:
echo   • Windows network interfaces sometimes don't support multicast properly
echo   • This is especially common with VPN, virtual, or tunneled interfaces
echo   • The warnings are harmless and don't affect P2P functionality
echo.
echo  Your node will still work correctly:
echo   • Peer discovery will work on supported interfaces
echo   • Local network discovery (same subnet) will function
echo   • Remote peer discovery via bootstrap nodes works fine
echo.
echo ️ To reduce warnings (optional):
echo   • Disable VPN software temporarily
echo   • Use Ethernet instead of WiFi
echo   • Or simply ignore them - they don't affect operation
echo.
echo  The warnings don't prevent:
echo   • Peer-to-peer connections
echo   • Service discovery and execution
echo   • Network communication
echo.
pause
