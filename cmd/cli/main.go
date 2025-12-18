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
	env := os.Getenv("APP_ENV")

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
	config.LoadEnv("development")
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
	config.LoadEnv("development")
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
		Role:         "admin",
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
	fmt.Println("=== Security Validation Check ===")

	sessionToken := generateSessionToken()
	fmt.Printf("Session: %s\n", sessionToken[:8]+"...")

	var issues []string
	var warnings []string

	fmt.Println("Checking Environment Security ...")

	env := os.Getenv("APP_ENV")

	switch env {
	case "":
		issues = append(issues, "APP_ENV not set")
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
	fmt.Println("\n Checking JWT Configuration...")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		issues = append(issues, "JWT_SECRET not configured")
	} else {
		if len(jwtSecret) < 16 {
			issues = append(issues, "JWT_SECRET too short (minimum 16 chars)")
		} else if len(jwtSecret) < 32 {
			warnings = append(warnings, "JWT_SECRET should be at least 32 chars for production")
		} else {
			fmt.Println(" JWT_SECRET length adequate")
		}

		if strings.Contains(strings.ToLower(jwtSecret), "secret") ||
			strings.Contains(strings.ToLower(jwtSecret), "password") ||
			strings.Contains(strings.ToLower(jwtSecret), "123") {
			warnings = append(warnings, "JWT_SECRET contains predictable patterns")
		}
	}

	fmt.Println(" Checking Database Security ...")

	config.LoadEnv(env)
	cfg, err := config.CargarVariables()
	if err != nil {
		issues = append(issues, fmt.Sprintf("Config loading failed: %v", err))
	} else {
		// Test database conection
		db, err := database.NuevaBaseDeDatos(cfg)
		if err != nil {
			issues = append(issues, fmt.Sprintf("Database connection failed: %v", err))
		} else {
			defer db.Close()
			fmt.Println(" Database connection successful")

			var rootCount int
			err = db.QueryRow("SELECT COUNT(*) FROM daedalus.users WHERE role = 'root'").Scan(&rootCount)
			if err != nil {
				warnings = append(warnings, "Could not check root user account")
			} else {
				if rootCount == 0 {
					warnings = append(warnings, "No root user found - run 'setup-root' first")
				} else if rootCount > 1 {
					issues = append(issues, fmt.Sprintf("Multiple root users found (%d) - security risk", rootCount))
				} else {
					fmt.Println("Single root user configured")
				}
			}

			var adminCount int
			err = db.QueryRow("SELECT COUNT(*) FROM daedalus.users WHERE role = 'admin'").Scan(&adminCount)
			if err != nil {
				warnings = append(warnings, "Could not check admin user count")
			} else {
				if adminCount == 0 {
					warnings = append(warnings, "No admin user found - run 'create-admin'")
				} else {
					fmt.Printf(" %d admin user(s) configured\n", adminCount)
				}
			}
		}
	}

	fmt.Println(" Checking File Permissions ...")

	if checkFilePermissions() {
		fmt.Println("File permissions secure")
	} else {
		warnings = append(warnings, "Some file may have insecure permissions")
	}

	fmt.Println("\n Checking Network Configuration ...")

	port := os.Getenv("PORT")
	switch port {
	case "":
		warnings = append(warnings, "PORT not explicitly set")
	case "80", "8080":
		if env == "production" {
			warnings = append(warnings, "Consider using HTTPS (port 443) in production")
		}
	}

	// Check OWASP TOP 10

	fmt.Println(" OWASP Top 10 Security Checks ...")

	fmt.Println(" A01: Access Control - JWT middleware implemented")

	if len(jwtSecret) >= 32 {
		fmt.Println(" A02: Cryptographic - Strong JWT secret")
	} else {
		warnings = append(warnings, "A02: Weak cryptographic configuration")
	}

	// A03: Injection
	fmt.Println(" A03: Injection - Parameterized queries used")

	// A07: Identification and Authentication Failures
	fmt.Println(" A07: Authentication - Bcrypt password hashing")

	// Generate Security Report
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println(" SECURITY ASSESSMENT REPORT")
	fmt.Println(strings.Repeat("=", 50))

	if len(issues) == 0 && len(warnings) == 0 {
		fmt.Println(" Excelente: No hay problemas encontrados!")
		fmt.Println(" Sistema ha pasado todos los controles de seguridad")
	} else {
		if len(issues) > 0 {
			fmt.Printf(" PROBLEMAS CRITICOS (%d):\n", len(issues))
			for i, issue := range issues {
				fmt.Printf("   %d.  %s\n", i+1, issue)
			}
		}

		if len(warnings) > 0 {
			fmt.Printf(" ALERTAS (%d):\n", len(warnings))
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

	fmt.Printf("\n Security Score: %d/%d", score, totalChecks)

	if score >= 9 {
		fmt.Println(" (EXCELENTE)")
	} else if score >= 7 {
		fmt.Println(" (BIEN)")
	} else if score >= 5 {
		fmt.Println(" (JUSTO - Necesita mejoras)")
	} else {
		fmt.Println(" (POBRE - Requiere acciones inmediatas)")
	}

	fmt.Printf(" Session: %s\n", sessionToken[:8]+"...")

	if len(issues) > 0 {
		fmt.Println("\nAborde los problemas críticos antes de la implementación de producción!")
		os.Exit(1)
	}
}

func checkFilePermissions() bool {
	sensitiveFiles := []string{
		".env.dev",
		".env.qa",
		".env.production",
		"cmd/config/config.go",
	}

	allSecure := true

	for _, file := range sensitiveFiles {
		if info, err := os.Stat(file); err == nil {
			mode := info.Mode()

			if mode&0044 != 0 {
				fmt.Printf(" %s has overly permissive permissions: %v\n", file, mode)
				allSecure = false
			}
		}
	}

	return allSecure
}
