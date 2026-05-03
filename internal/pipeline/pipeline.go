package pipeline

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type Step struct {
	Name    string
	Command string
}

func RunPipeline(repo string) {
	log.Println("🚀 Pipeline started for:", repo)

	// 🔥 version unique
	version := fmt.Sprintf("v%d", time.Now().Unix())

	// 🔐 token GitHub (env variable)
	token := os.Getenv("GITHUB_TOKEN")

	if token == "" {
		log.Println("❌ GITHUB_TOKEN not set")
		return
	}

	// 🐳 image Docker versionnée
	imageName := "nourhenhachem/innodeploy-app:" + version

	// 🔗 remote Git avec auth
	remote := fmt.Sprintf(
		"https://hachemn:%s@github.com/hachemn/innodeploy-pfe.git",
		token,
	)

	steps := []Step{
		{"Cleanup", "rm -rf project"},
		{"Clone", "git clone " + repo + " project"},

		// 🐳 Build
		{"Build Docker Image", "cd project && docker build -t innodeploy-app ."},

		// 🏷️ Tag version
		{"Tag Image", "docker tag innodeploy-app " + imageName},

		// 🚀 Push Docker
		{"Push Image", "docker push " + imageName},

		// 🧠 GitOps → update YAML
		{"Update YAML", fmt.Sprintf(
			"sed -i 's|image:.*|image: %s|' k8s/deployment.yaml",
			imageName,
		)},

		// 🔐 config remote avec token
		{"Set Git Remote", "git remote set-url origin " + remote},

		// 📝 commit
		{"Git Commit", fmt.Sprintf(
			"git add . && git commit -m 'deploy %s'",
			version,
		)},

		// 📤 push
		{"Git Push", "git push"},
	}

	for _, step := range steps {
		runStep(step)
	}

	log.Println("✅ Pipeline finished 🚀 version:", version)
}

func runStep(step Step) {
	log.Println("➡️ Step:", step.Name)

	command := exec.Command("sh", "-c", step.Command)
	output, err := command.CombinedOutput()

	if err != nil {
		log.Println("❌ Error in step:", step.Name)
		log.Println(string(output))
		return
	}

	log.Println("✅ Success:", step.Name)
	log.Println(string(output))
}