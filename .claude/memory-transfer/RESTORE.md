# Memory Transfer Instructions (Windows)

These files are Claude Code project memory exported from the macOS laptop on 2026-04-20.
To restore them on your Windows PC after cloning the repo:

## Steps

1. Clone the repo (e.g. `gh repo clone fajarmf10/secret-scraper`) and cd into it.
2. Open Claude Code in the repo once — this creates its project folder under
   `%USERPROFILE%\.claude\projects\<slug>\`. The slug is the Windows path with
   slashes replaced by dashes, e.g.
   `-C--Users-fajar-projects-secret-scraper`.
3. Copy memories into the Windows slot:
   ```powershell
   $memDir = "$env:USERPROFILE\.claude\projects\<slug>\memory"
   New-Item -ItemType Directory -Force -Path $memDir | Out-Null
   Copy-Item .claude\memory-transfer\*.md $memDir\ -Exclude RESTORE.md,PROJECT_MEMORY_INDEX.md
   # Rename the index back to MEMORY.md inside the memory dir:
   Copy-Item .claude\memory-transfer\PROJECT_MEMORY_INDEX.md $memDir\MEMORY.md -ErrorAction SilentlyContinue
   ```
4. Or, in Claude Code, just say: "Read `.claude/memory-transfer/` and re-save each as project memory."
