package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/sha1"
    "crypto/sha512"
    "encoding/base64"
    "encoding/hex"
    "encoding/json"
    "bytes"
    "fmt"
    "log"
    "net/http"
    "io"
    "os"
    "runtime"
    "time"
    "os/exec"
)

const (
    verifyURL = "https://yourdomain.tld/api/auth/verify/"
    releaseURL = "https://github.com/hashburst/HashburstMiner/releases/tag/HashburstMiner"
    secretKey = "Chiper secret key"
    secretIV  = "IV initialization vector"
)

type User struct {
    Email         string    `json:"email"`
    APIKey        string    `json:"apikey"`
    ReferralCode  string    `json:"referral_code"`
    DogeWallet    string    `json:"doge_wallet"`
    BTCWallet     string    `json:"btc_wallet"`
    LicenseExpiry time.Time `json:"license_expiry"`
}

var users = make(map[string]User)

type VerificationResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}

func saveUsersToFile() {
    file, err := os.Create("users.json")
    if err != nil {
        log.Fatal("Errore nel creare il file:", err)
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    err = encoder.Encode(users)
    if err != nil {
        log.Fatal("Encoding error:", err)
    }
}

func loadUsersFromFile() {
    file, err := os.Open("users.json")
    if err != nil {
        if os.IsNotExist(err) {
            return
        }
        log.Fatal("Opening file error:", err)
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    err = decoder.Decode(&users)
    if err != nil {
        log.Fatal("Decoding error:", err)
    }
}

func verifyLicense(w http.ResponseWriter, r *http.Request) {
    var reqUser User
    err := json.NewDecoder(r.Body).Decode(&reqUser)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    user, exists := users[reqUser.Email]
    if !exists {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    if !isLicenseValid(user.LicenseExpiry) {
        http.Error(w, "License expired", http.StatusForbidden)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "message":       "License valid",
        "license_expiry": user.LicenseExpiry.Format(time.RFC3339),
    })
}

func isLicenseValid(licenseExpiry time.Time) bool {
    return time.Now().Before(licenseExpiry)
}

// Encrypt and decrypt function
func encryptDecrypt(input, key, iv string, action string) string {
    keyHash := hashSHA512(key)[:32] // First 32 bytes of hash for key
    ivHash := hashSHA512(iv)[:16]   // Firsr 16 bytes of hash per initialization vector

    block, err := aes.NewCipher([]byte(keyHash))
    if err != nil {
        log.Fatal("Errore nel creare il cifrario:", err)
    }

    if action == "encrypt" {
        cipherText := make([]byte, len(input))
        mode := cipher.NewCBCEncrypter(block, []byte(ivHash))
        mode.CryptBlocks(cipherText, []byte(input))
        return base64.StdEncoding.EncodeToString(cipherText)
    } else if action == "decrypt" {
        decoded, _ := base64.StdEncoding.DecodeString(input)
        decrypted := make([]byte, len(decoded))
        mode := cipher.NewCBCDecrypter(block, []byte(ivHash))
        mode.CryptBlocks(decrypted, decoded)
        return string(decrypted)
    }

    return ""
}

func hashSHA512(input string) string {
    hash := sha512.New()
    hash.Write([]byte(input))
    return hex.EncodeToString(hash.Sum(nil))
}

func generateProofOfHistory(data string) map[string]string {
    timestamp := time.Now().Unix()
    combinedData := fmt.Sprintf("%s%d", data, timestamp)
    proof := hashSHA512(combinedData)
    return map[string]string{
        "data":      data,
        "timestamp": fmt.Sprintf("%d", timestamp),
        "proof":     proof,
    }
}

func registerUser(w http.ResponseWriter, r *http.Request) {
    var newUser User
    err := json.NewDecoder(r.Body).Decode(&newUser)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

   if !verifyWithHashburst(newUser.Email, newUser.APIKey, newUser.ReferralCode) {
        http.Error(w, "Verification failed with Hashburst", http.StatusUnauthorized)
        return
   }

   newUser.LicenseExpiry = time.Now().AddDate(0, 6, 0)

   users[newUser.Email] = newUser
   saveUsersToFile()

   json.NewEncoder(w).Encode(map[string]string{
        "apikey":        newUser.APIKey,
        "license_expiry": newUser.LicenseExpiry.Format(time.RFC3339),
   })
}

// Function POST request using cURL
func verifyWithHashburst(email, apikey, referralCode string) bool {
    reqBody := map[string]string{
        "email":        email,
        "apikey":       apikey,
        "referral_code": referralCode,
    }

    jsonReq, err := json.Marshal(reqBody)
    if err != nil {
        log.Println("Error marshalling request:", err)
        return false
    }

    // cURL command for per POST request
    curlCmd := exec.Command("curl", "-X", "POST", verifyURL, "-H", "Content-Type: application/json", "-d", string(jsonReq))

    // Command cURL execution
    curlOutput, err := curlCmd.Output()
    if err != nil {
        log.Println("Error executing cURL command:", err)
        return false
    }

    var result map[string]interface{}
    err = json.Unmarshal(curlOutput, &result)
    if err != nil {
        //log.Println("Error unmarshalling response:", err)
        //return false
    }

    return result["status"] == "success"
}

// Test network speed
func testNetworkSpeed() error {
    var cmd *exec.Cmd
    if runtime.GOOS == "windows" {
        cmd = exec.Command("cmd", "/C", "ping", "8.8.8.8", "-n", "10")
    } else {
        cmd = exec.Command("ping", "-c", "10", "8.8.8.8")
    }

    var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        return fmt.Errorf("error during network speed test: %v", err)
    }

    output := out.String()
    fmt.Println("Network speed test output:\n", output)
    return nil
}

// Function for downloading Miner last version for free
func downloadMiner() error {
    var fileName string
    switch runtime.GOOS {
    case "windows":
        fileName = "HashburstMiner.exe"
    case "darwin":
        if runtime.GOARCH == "amd64" {
            fileName = "HashburstMiner"
        } else if runtime.GOARCH == "arm64" {
            fileName = "HashburstMiner-Mac-arm64"
            os.Rename("HashburstMiner-Mac-arm64", "HashburstMiner")
        }
    case "linux":
        fileName = "HashburstMiner-Linux"
        os.Rename("HashburstMiner-Linux", "HashburstMiner")
    }

    resp, err := http.Get(releaseURL)
    if err != nil {
        return fmt.Errorf("Errore last release: %v", err)
    }
    defer resp.Body.Close()

    out, err := os.Create(fileName)
    if err != nil {
        return fmt.Errorf("Error file creation: %v", err)
    }
    defer out.Close()

    _, err = io.Copy(out, resp.Body)
    if err != nil {
        return fmt.Errorf("Error file saving: %v", err)
    }

    log.Println("Completed download", fileName)
    return nil
}

// Function to run hashburstMiner
func startMiner() error {
    err := downloadMiner()
    if err != nil {
        return fmt.Errorf("Download miner error: %v", err)
    }

    var cmd *exec.Cmd
    switch runtime.GOOS {
    case "windows":
        cmd = exec.Command("./HashburstMiner.exe", "-o", "cryptonightr.auto.nicehash.com:9200", "-u", "Mining Wallet Nicehash", "-k", "--nicehash", "--coin", "btc", "-a", "cn/r", "--rig-id=freeNodes", "-l")
    case "darwin", "linux":
        cmd = exec.Command("./HashburstMiner", "-o", "cryptonightr.auto.nicehash.com:9200", "-u", "Mining Wallet Nicehash", "-k", "--nicehash", "--coin", "btc", "-a", "cn/r", "--rig-id=freeNodes", "-l")
    }

    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    err = cmd.Start()
    if err != nil {
        return fmt.Errorf("Start mining error: %v", err)
    }

    log.Println("Miner works!")
    return cmd.Wait()
}

func main() {
    loadUsersFromFile()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./templates/index.html")
    })

    http.HandleFunc("/register", registerUser)
    http.HandleFunc("/verify", verifyLicense)

    testNetworkSpeed()

    downloadMiner()

    startMiner()

    log.Println("Starting server on :8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))    
}
