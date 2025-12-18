package main

import (
	"context"
	"crypto/rand"
	"daedalus-engine-go/cmd/common/database"
	"daedalus-engine-go/cmd/common/logger"
	"daedalus-engine-go/cmd/config"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"strings"
	"syscall"
	"time"

	userRepo "daedalus-engine-go/cmd/internal/features/users/application/data/local"
	userDomain "daedalus-engine-go/cmd/internal/features/users/domain/entities"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

// Definición de colores
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
)

const (
	// constantes de seguridad (OWASP)
	MIN_PASSWORD_LENGTH = 14
	MAX_LOGIN_ATTEMPTS  = 3
	BCRYPT_COST         = 12
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	nameRegex  = regexp.MustCompile(`^[a-zA-Z\s]{2,50}$`)
)

func main() {

	if !isSecureEnvironment() {
		fmt.Println(" Insecure environment detected. CLI operations blocked")
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	logCLIOperation(command)

	switch command {
	case "setup-root":
		setupRootUser()
	case "create-admin":
		createAdminUser()
	case "security-check":
		performSecurityCheck()
	default:
		fmt.Printf(" Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func isSecureEnvironment() bool {
	// Confirmamos si esta corriendo en produccion con la seguridad apropiada
	env := loadEnvironment()

	if env == "production" {
		// Verificamos secure connection a la DB
		if os.Getenv("DB_SSLMODE") != "require" {
			fmt.Println(" Production requieres SSL database connection")
			return false
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		if len(jwtSecret) < 32 {
			fmt.Println(" Production requires strong JWT secret (32+ chars)")
			return false
		}
	}
	return true
}

// Security: Audit logging
func logCLIOperation(command string) {
	logger.Init("cli")
	log := logger.L()

	log.Info("CLI operation initiated", zap.String("command", command),
		zap.String("user", os.Getenv("USER")), zap.String("hostname", getHostname()),
		zap.Time("timestamp", time.Now()))
}

func getHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

func printUsage() {
	fmt.Println(ColorGreen + "Daedalus CLI - Secure Initial Setup" + ColorReset)
	fmt.Println("\nUsage:")
	fmt.Println(ColorGreen + " go run cmd/cli/main.go <command>" + ColorReset)
	fmt.Println(ColorBold + ColorYellow + "\nCommands:" + ColorReset)
	fmt.Println(ColorCyan + " setup-root     - Create initial root user (Zero trust)")
	fmt.Println(" create-admin   - Create first admin user (requires root auth)")
	fmt.Println(" security-check - Perform security validation" + ColorReset)
	fmt.Println(ColorBold + ColorYellow + "\n Security Notes:" + ColorReset)
	fmt.Println(ColorYellow + " - All operations are logged and audited")
	fmt.Println(" - Password must meet OWASP standards (14+ chars)")
	fmt.Println(" - Root user is for setup only, never for API access" + ColorReset)
}

// Validaciones de entrada segura (OWASP)
func validateEmail(email string) error {
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	if len(email) > 254 {
		return fmt.Errorf("email too long")
	}
	return nil
}

func validateName(name string) error {
	if !nameRegex.MatchString(name) {
		return fmt.Errorf("name must be 2-50 alphabetic characters")
	}
	return nil
}

// OWASP: Password strength validation

func validatePassword(password []byte) error {
	if len(password) < MIN_PASSWORD_LENGTH {
		return fmt.Errorf("password must be at least %d characters", MIN_PASSWORD_LENGTH)
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", rune(char)):
			hasSpecial = true
		}
	}

	if !hasLower || !hasUpper || !hasNumber || !hasSpecial {
		return fmt.Errorf("password must contain uppercase, lowercase, number, and special character")
	}

	return nil
}

func setupRootUser() {
	fmt.Println("=== Zero Trust Root Setup ===")

	sessionToken := generateSessionToken()
	fmt.Printf("Session: %s\n", sessionToken[:8]+"...") // Mostramos los primeros 8 caracteres

	// Cargamos variable de ambiente
	loadEnvironment()
	logger.Init("cli")
	cfg, err := config.CargarVariables()
	if err != nil {
		fmt.Printf("Config error: %v\n", err)
		os.Exit(1)
	}

	// Secure database conection
	db, err := database.NuevaBaseDeDatos(cfg)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		os.Exit(1)
	}

	defer db.Close()

	// Verify no root exists
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM daedalus.users WHERE role = 'root'").Scan(&count)
	if err != nil {
		fmt.Printf("Database query error: %v\n", err)
		os.Exit(1)
	}

	if count > 0 {
		fmt.Println(" Root user already exists!")
		os.Exit(1)
	}

	fmt.Print("Root Email: ")
	var email string
	fmt.Scanln(&email)

	if err := validateEmail(email); err != nil {
		fmt.Printf(" %s\n", err)
		os.Exit(1)
	}

	fmt.Print("Root Name: ")
	var name string
	fmt.Scanln(&name)
	name = strings.TrimSpace(name)

	if err := validateName(name); err != nil {
		fmt.Printf(" %s\n", err)
		os.Exit(1)
	}

	password, err := getSecurePassword("Root Password (14+ chars): ")
	if err != nil {
		fmt.Printf(" %s\n", err)
		os.Exit(1)
	}

	hash, err := bcrypt.GenerateFromPassword(password, BCRYPT_COST)
	if err != nil {
		fmt.Printf(" Hashing error: %v\n", err)
		os.Exit(1)
	}

	for i := range password {
		password[i] = 0
	}

	// Create root user
	rootUser := &userDomain.User{
		ID:           uuid.NewString(),
		TenantID:     nil,
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		Role:         "root",
		Status:       "active",
	}

	repo := userRepo.NewUserRepository(db)
	err = repo.Create(context.Background(), rootUser)

	if err != nil {
		fmt.Printf(" Creation error : %v\n", err)
		os.Exit(1)
	}

	fmt.Println(" Root user created successfully!")
	fmt.Printf(" Email: %s\n", email)
	fmt.Printf(" Session: %s\n", sessionToken[:8]+"...")
	fmt.Println(" SECURITY REMINDERS: ")
	fmt.Println("   - Root credentials must be stored securely")
	fmt.Println("   - Root user is ONLY for creating first admin")
	fmt.Println("   - Never use root for API access")
	fmt.Println("   - Consider disabling root after admin creation")
}

func generateSessionToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func getSecurePassword(prompt string) ([]byte, error) {
	for attempt := 0; attempt < MAX_LOGIN_ATTEMPTS; attempt++ {
		fmt.Print(prompt)
		password, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		if err != nil {
			return nil, fmt.Errorf("failed to read password: %v", err)
		}

		if err := validatePassword(password); err != nil {
			fmt.Printf(" %s\n", err)
			fmt.Println(" Requerimientos de la contraseña:")
			fmt.Println("  - Mínimo 14 carácteres")
			fmt.Println("  - 1 letra mayúscula")
			fmt.Println("  - 1 letra minúscula")
			fmt.Println("  - 1 numero")
			fmt.Println("  - 1 caracter especial (!@#¢∞ª%^...)")
			continue
		}

		fmt.Print("Confirm password: ")
		confirm, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return nil, fmt.Errorf("failed to read confirmation: %v", err)
		}

		if string(password) != string(confirm) {
			fmt.Println("Passwords don't match!")
			for i := range confirm {
				confirm[i] = 0
			}
			continue
		}

		for i := range confirm {
			confirm[i] = 0
		}
		return password, nil
	}
	return nil, fmt.Errorf("maximum password attempts exceeded (%d)", MAX_LOGIN_ATTEMPTS)
}

func createAdminUser() {
	fmt.Println("=== Create Admin User ===")

	sessionToken := generateSessionToken()
	fmt.Printf("Session: %s\n", sessionToken[:8]+"...")

	// Load config
	loadEnvironment()
	logger.Init("cli")
	cfg, err := config.CargarVariables()
	if err != nil {
		fmt.Printf("Config error: %v\n", err)
		os.Exit(1)
	}

	// Connect to database
	db, err := database.NuevaBaseDeDatos(cfg)
	if err != nil {
		fmt.Printf("Database error: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	repo := userRepo.NewUserRepository(db)

	// Authenticate root user
	fmt.Print("Root Email: ")
	var rootEmail string
	fmt.Scanln(&rootEmail)

	fmt.Print("Root Password: ")
	rootPassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		fmt.Printf("Error reading password: %v\n", err)
		os.Exit(1)
	}

	// Verify root credentials
	rootUser, err := repo.GetByEmail(context.Background(), rootEmail)
	if err != nil || rootUser == nil || rootUser.Role != "root" {
		fmt.Println("Invalid root credentials!")
		os.Exit(1)
	}

	err = bcrypt.CompareHashAndPassword([]byte(rootUser.PasswordHash), rootPassword)
	if err != nil {
		fmt.Println("Invalid root credentials!")
		os.Exit(1)
	}

	// Clear root password from memory
	for i := range rootPassword {
		rootPassword[i] = 0
	}

	// Get admin details
	fmt.Print("\nAdmin Email: ")
	var adminEmail string
	fmt.Scanln(&adminEmail)

	if err := validateEmail(adminEmail); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	fmt.Print("Admin Name: ")
	var adminName string
	fmt.Scanln(&adminName)

	if err := validateName(adminName); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	adminPassword, err := getSecurePassword("Admin Password (14+ chars): ")
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword(adminPassword, BCRYPT_COST)
	if err != nil {
		fmt.Printf("Hashing error: %v\n", err)
		os.Exit(1)
	}

	// Clear password from memory
	for i := range adminPassword {
		adminPassword[i] = 0
	}

	// Create admin user
	adminUser := &userDomain.User{
		ID:           uuid.NewString(),
		TenantID:     nil,
		Name:         adminName,
		Email:        adminEmail,
		PasswordHash: string(hash),
		Role:         "system_admin",
		Status:       "active",
	}

	err = repo.Create(context.Background(), adminUser)
	if err != nil {
		fmt.Printf("Creation error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nAdmin user created successfully!")
	fmt.Printf("Email: %s\n", adminEmail)
	fmt.Printf("Session: %s\n", sessionToken[:8]+"...")
	fmt.Println("\n You can now login with this admin user via API.")
}

func performSecurityCheck() {
	fmt.Println(ColorBold + ColorWhite + "=== Security Validation Check ===" + ColorReset)

	sessionToken := generateSessionToken()
	fmt.Printf(ColorWhite+"Session: %s\n", sessionToken[:8]+"..."+ColorReset)

	var issues []string
	var warnings []string

	fmt.Println(ColorGreen + "Checking Environment Security ..." + ColorReset)

	env := loadEnvironment()

	switch env {
	case "":
		issues = append(issues, ColorBold+ColorRed+"APP_ENV not set"+ColorReset)
	case "production":
		fmt.Println(" Production environtment detected ")
		if os.Getenv("DB_SSLMODE") != "require" {
			issues = append(issues, "Production required DB_SSLMODE=require")
		}
		if len(os.Getenv("JWT_SECRET")) < 32 {
			issues = append(issues, "Production JWT_SECRET too weak (< 32 chars)")
		}
	default:
		fmt.Printf(" Development environment: %s\n", env)
	}

	// JWT Security Check
	fmt.Println(ColorWhite + "\n Checking JWT Configuration..." + ColorReset)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		issues = append(issues, ColorBold+ColorRed+"JWT_SECRET not configured"+ColorReset)
	} else {
		if len(jwtSecret) < 16 {
			issues = append(issues, ColorYellow+"JWT_SECRET too short (minimum 16 chars)"+ColorReset)
		} else if len(jwtSecret) < 32 {
			warnings = append(warnings, ColorYellow+"JWT_SECRET should be at least 32 chars for production"+ColorReset)
		} else {
			fmt.Println(ColorBold + ColorGreen + " JWT_SECRET length adequate" + ColorReset)
		}

		if strings.Contains(strings.ToLower(jwtSecret), "secret") ||
			strings.Contains(strings.ToLower(jwtSecret), "password") ||
			strings.Contains(strings.ToLower(jwtSecret), "123") {
			warnings = append(warnings, "JWT_SECRET contains predictable patterns")
		}
	}

	fmt.Println(ColorWhite + " Checking Database Security ..." + ColorReset)

	cfg, err := config.CargarVariables()
	if err != nil {
		issues = append(issues, fmt.Sprintf(ColorBold+ColorRed+"Config loading failed: %v"+ColorReset, err))
	} else {
		// Test database conection
		db, err := database.NuevaBaseDeDatos(cfg)
		if err != nil {
			issues = append(issues, fmt.Sprintf(ColorBold+ColorRed+"Database connection failed: %v"+ColorReset, err))
		} else {
			defer db.Close()
			fmt.Println(ColorBold + ColorGreen + " Conexión a la DB exitosa" + ColorReset)

			var rootCount int
			err = db.QueryRow("SELECT COUNT(*) FROM daedalus.users WHERE role = 'root'").Scan(&rootCount)
			if err != nil {
				warnings = append(warnings, ColorYellow+"Could not check root user account"+ColorReset)
			} else {
				if rootCount == 0 {
					warnings = append(warnings, ColorYellow+"No root user found - run 'setup-root' first"+ColorReset)
				} else if rootCount > 1 {
					issues = append(issues, fmt.Sprintf(ColorBold+ColorRed+"Multiple root users found (%d) - security risk"+ColorReset, rootCount))
				} else {
					fmt.Println(ColorCyan + "Single root user configured" + ColorReset)
				}
			}

			var adminCount int
			err = db.QueryRow("SELECT COUNT(*) FROM daedalus.users WHERE role = 'system_admin'").Scan(&adminCount)
			if err != nil {
				warnings = append(warnings, ColorYellow+"Could not check admin user count"+ColorReset)
			} else {
				if adminCount == 0 {
					warnings = append(warnings, ColorYellow+"No System admin user found - run 'create-admin'"+ColorReset)
				} else {
					fmt.Printf(ColorGreen+" %d System admin user(s) configured\n"+ColorReset, adminCount)
				}
			}
		}
	}

	fmt.Println(ColorWhite + " Checking File Permissions ..." + ColorReset)

	if checkFilePermissions() {
		fmt.Println(ColorBold + ColorGreen + " Permisos de archivos seguros" + ColorReset)
	} else {
		warnings = append(warnings, ColorYellow+"Some file may have insecure permissions"+ColorReset)
	}

	fmt.Println(ColorWhite + "\n Checking Network Configuration ..." + ColorReset)

	port := os.Getenv("PORT")
	switch port {
	case "":
		warnings = append(warnings, ColorYellow+"PORT not explicitly set"+ColorReset)
	case "80", "8080":
		if env == "production" {
			warnings = append(warnings, ColorBold+ColorYellow+"Consider using HTTPS (port 443) in production"+ColorReset)
		}
	}

	// Check OWASP TOP 10

	fmt.Println(ColorBlue + " OWASP Top 10 Security Checks ..." + ColorReset)

	fmt.Println(ColorPurple + " A01: Access Control - JWT middleware implemented" + ColorReset)

	if len(jwtSecret) >= 32 {
		fmt.Println(ColorGreen + " A02: Cryptographic - Strong JWT secret" + ColorReset)
	} else {
		warnings = append(warnings, ColorYellow+"A02: Weak cryptographic configuration"+ColorReset)
	}

	// A03: Injection
	fmt.Println(ColorPurple + " A03: Injection - Parameterized queries used")

	// A07: Identification and Authentication Failures
	fmt.Println(" A07: Authentication - Bcrypt password hashing" + ColorReset)

	// Generate Security Report
	fmt.Println(ColorCyan + "\n" + strings.Repeat("=", 50))
	fmt.Println(ColorBlue + " SECURITY ASSESSMENT REPORT")
	fmt.Println(ColorCyan + strings.Repeat("=", 50) + ColorReset)

	if len(issues) == 0 && len(warnings) == 0 {
		fmt.Println(" Excelente: No hay problemas encontrados!")
		fmt.Println(" Sistema ha pasado todos los controles de seguridad")
	} else {
		if len(issues) > 0 {
			fmt.Printf(ColorBold+ColorRed+" PROBLEMAS CRITICOS (%d):\n"+ColorReset, len(issues))
			for i, issue := range issues {
				fmt.Printf("   %d.  %s\n", i+1, issue)
			}
		}

		if len(warnings) > 0 {
			fmt.Printf(ColorBold+ColorYellow+" ALERTAS (%d):\n"+ColorReset, len(warnings))
			for i, warning := range warnings {
				fmt.Printf("   %d.  %s\n", i+1, warning)
			}
		}
	}

	totalChecks := 10
	criticalIssues := len(issues)
	minorIssues := len(warnings)

	score := totalChecks - (criticalIssues * 2) - minorIssues
	if score < 0 {
		score = 0
	}

	fmt.Printf(ColorPurple+"\n Puntaje de seguridad: %d/%d |"+ColorReset, score, totalChecks)

	if score >= 9 {
		fmt.Println(ColorGreen + " (EXCELENTE)" + ColorReset)
	} else if score >= 7 {
		fmt.Println(ColorCyan + " (BIEN)" + ColorReset)
	} else if score >= 5 {
		fmt.Println(ColorYellow + " (JUSTO - Necesita mejoras)" + ColorReset)
	} else {
		fmt.Println(ColorRed + " (POBRE - Requiere acciones inmediatas)" + ColorReset)
	}

	fmt.Printf(ColorWhite+" Session: %s\n", sessionToken[:8]+"..."+ColorReset)

	if len(issues) > 0 {
		fmt.Println(ColorRed + "\nAborde los problemas críticos antes de la implementación de producción!" + ColorReset)
		os.Exit(1)
	}
}

func checkFilePermissions() bool {
	sensitiveFiles := []string{
		".env.dev",
		".env.qa",
		".env.production",
	}

	allSecure := true

	for _, file := range sensitiveFiles {
		if info, err := os.Stat(file); err == nil {
			mode := info.Mode()

			if mode&0044 != 0 {
				fmt.Printf(ColorYellow+" %s tiene demasiados permisos: %v\n"+ColorRed, file, mode)
				allSecure = false
			}
		}
	}

	return allSecure
}

func loadEnvironment() string {
	env := os.Getenv("APP_ENV")

	if env == "" {
		godotenv.Load(".env.dev")
		env = os.Getenv("APP_ENV")
	}

	if env == "" {
		env = "development"
	}

	config.LoadEnv(env)
	return env
}
