


# Distributed Model Training System

## Overview

The Distributed Model Training System is designed to efficiently utilize underutilized computing resources, particularly GPUs, on a remote system. This project enables the remote system to execute training tasks by receiving commands from a sender system. The sender system provides the necessary data and training scripts, which are then executed on the remote system. The results of the training are committed back to the GitHub repository from which the code and data were sourced. This allows for seamless collaboration and version control, ensuring that all changes are tracked and easily accessible.

## Project Components

### 1. **Executor**

The `executor.go` file serves as the core component that listens for incoming commands. It initiates a tunnel using Ngrok to expose the local server to the public internet, allowing it to receive webhook requests. The executor fetches the public URL from Ngrok and updates the configuration file (`config.json`) with the current public URL and timestamp. The executor also handles the execution of received training commands.

### 2. **Sender**

The `sender.go` file is responsible for sending instructions to the executor. It reads the public URL from the `config.json` file and formulates the necessary commands to execute the training script. The sender system provides the data URL and training command, which are then sent to the executor via an HTTP POST request.

### 3. **Configuration Management**

The project uses a JSON configuration file (`config.json`) to store important information such as the public URL, IP address, port, and timestamp. This configuration file is automatically updated by the executor to ensure that it contains the latest information.

## How It Works

1. **Setup**: 
   - The executor starts and initializes Ngrok to create a public URL for the local server.
   - The executor updates the `config.json` file with the public URL.

2. **Receiving Instructions**:
   - The sender monitors changes in the GitHub repository (using webhooks).
   - Upon detecting a change, the sender retrieves the public URL from the `config.json` file and sends training instructions to the executor.

3. **Executing Training**:
   - The executor receives the training command and executes it, utilizing the GPU resources of the remote system.

4. **Committing Results**:
   - Once the training is complete, the executor commits the results back to the GitHub repository.

## Getting Started

### Prerequisites

- Go installed on your system
- Ngrok for tunneling
- Access to a GitHub repository for code and data storage

### Installation

1. Clone the repository to your local machine:
   <<<bash
   git clone https://github.com/yourusername/yourrepository.git
   cd yourrepository
   <<<

2. Install Ngrok and set up your authtoken in the `ngrok.yml` configuration file.

3. Build the Go programs:
   <<<bash
   go build executor.go
   go build sender.go
   <<<

4. Run the executor:
   <<<bash
   ./executor
   <<<

5. Start the sender when ready to send instructions.

## Example Usage

1. Ensure that the executor is running and has started Ngrok.
2. Update your GitHub repository with the training script and data URL.
3. Push changes to the main branch of the repository to trigger the sender.
4. The sender will automatically send the training instructions to the executor.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

For any inquiries, please reach out to [your.email@example.com].
<<<
