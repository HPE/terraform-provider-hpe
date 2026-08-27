# Test suite for install-hpe-provider-windows.ps1
# Run with: Invoke-Pester -Path .\install-hpe-provider-windows.Tests.ps1

BeforeAll {
    # Store original location and environment
    $script:OriginalLocation = Get-Location
    $script:ScriptPath = Join-Path $PSScriptRoot ".." "install-hpe-provider-windows.ps1"
}

Describe "install-hpe-provider-windows.ps1" {
    BeforeEach {
        # Create temporary test directory
        $script:TestAppData = Join-Path $env:TEMP "terraform-test-$(New-Guid)"
        New-Item -ItemType Directory -Path $script:TestAppData -Force | Out-Null
        
        # Override AppData for testing
        $env:APPDATA = $script:TestAppData
    }
    
    AfterEach {
        # Clean up test directory
        Set-Location $script:OriginalLocation
        if (Test-Path $script:TestAppData) {
            Remove-Item -Path $script:TestAppData -Recurse -Force -ErrorAction SilentlyContinue
        }
    }
    
    Context "Installation with specific version" {
        It "creates correct directory structure" {
            & $script:ScriptPath -VERSION "v0.1.0"
            
            $expectedPath = Join-Path $env:APPDATA "terraform.d\plugins\registry.terraform.io\hpe\hpe\0.1.0\windows_amd64"
            Test-Path $expectedPath | Should -Be $true
        }
        
        It "downloads and extracts provider binary" {
            & $script:ScriptPath -VERSION "v0.1.0"
            
            $binaryPath = Join-Path $env:APPDATA "terraform.d\plugins\registry.terraform.io\hpe\hpe\0.1.0\windows_amd64\terraform-provider-hpe_v0.1.0.exe"
            Test-Path $binaryPath | Should -Be $true
        }
        
        It "removes zip file after extraction" {
            & $script:ScriptPath -VERSION "v0.1.0"
            
            $zipFiles = Get-ChildItem -Path (Join-Path $env:APPDATA "terraform.d") -Filter "*.zip" -Recurse
            $zipFiles.Count | Should -Be 0
        }
        
        It "removes temporary extraction folder" {
            & $script:ScriptPath -VERSION "v0.1.0"
            
            $extractPath = Join-Path $env:APPDATA "terraform.d\plugins\registry.terraform.io\hpe\hpe\0.1.0\windows_amd64"
            $tempFolders = Get-ChildItem -Path $extractPath -Directory -Filter "terraform-provider-hpe_*"
            $tempFolders.Count | Should -Be 0
        }
    }
    
    Context "Installation without explicit version" {
        It "downloads latest version when VERSION not specified" {
            & $script:ScriptPath
            
            $hpePath = Join-Path $env:APPDATA "terraform.d\plugins\registry.terraform.io\hpe\hpe"
            Test-Path $hpePath | Should -Be $true
            
            # At least one version directory should exist
            $versionDirs = Get-ChildItem -Path $hpePath -Directory
            $versionDirs.Count | Should -BeGreaterThan 0
        }
    }
    
    Context "Reinstallation" {
        It "can reinstall over existing installation" {
            # Install once
            & $script:ScriptPath -VERSION "v0.1.0"
            $firstInstall = Test-Path (Join-Path $env:APPDATA "terraform.d\plugins\registry.terraform.io\hpe\hpe\0.1.0\windows_amd64\terraform-provider-hpe_v0.1.0.exe")
            
            # Install again
            & $script:ScriptPath -VERSION "v0.1.0"
            $secondInstall = Test-Path (Join-Path $env:APPDATA "terraform.d\plugins\registry.terraform.io\hpe\hpe\0.1.0\windows_amd64\terraform-provider-hpe_v0.1.0.exe")
            
            $firstInstall | Should -Be $true
            $secondInstall | Should -Be $true
        }
    }
    
    Context "Error handling" {
        It "handles invalid version gracefully" {
            & $script:ScriptPath -VERSION "v999.999.999"
            
            # Should not create version directory for invalid version
            $invalidPath = Join-Path $env:APPDATA "terraform.d\plugins\registry.terraform.io\hpe\hpe\999.999.999"
            Test-Path $invalidPath | Should -Be $false
        }
        
        It "cleans up on download failure" {
            & $script:ScriptPath -VERSION "v999.999.999"
            
            # Verify cleanup occurred
            $basePath = Join-Path $env:APPDATA "terraform.d\plugins\registry.terraform.io\hpe\hpe"
            if (Test-Path $basePath) {
                $invalidDirs = Get-ChildItem -Path $basePath -Directory -Filter "999.999.999"
                $invalidDirs.Count | Should -Be 0
            }
        }
    }
    
    Context "Return to original location" {
        It "returns to user's original directory after installation" {
            $beforeLocation = Get-Location
            & $script:ScriptPath -VERSION "v0.1.0"
            $afterLocation = Get-Location
            
            # Script should return to original location (handled by our cleanup)
            Set-Location $beforeLocation
            $afterLocation.Path | Should -Not -Match "terraform.d"
        }
    }
}
