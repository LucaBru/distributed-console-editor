#!/usr/bin/env python3
import subprocess
import time
import argparse
import os
import signal
import sys
from concurrent.futures import ThreadPoolExecutor

def run_client(client_id, command, duration):
    """Run a single client process for the specified duration."""
    print(f"Starting client {client_id}")
    
    # Create a log file for this client
    log_file = open(f"analysis_log/client_{client_id}.log", "w")
    
    # Start the process
    process = subprocess.Popen(
        command,
        shell=True,
        stdout=log_file,
        stderr=subprocess.STDOUT,
        text=True
    )
    
    # Wait for the specified duration
    time.sleep(duration)
    
    # Terminate the process
    try:
        process.terminate()
        process.wait(timeout=5)  # Give it 5 seconds to terminate gracefully
    except subprocess.TimeoutExpired:
        process.kill()  # Force kill if it doesn't terminate
    
    log_file.close()
    print(f"Client {client_id} finished")
    return client_id

def main():
    parser = argparse.ArgumentParser(description="Launch multiple clients for distributed text editor testing")
    parser.add_argument("--clients", type=int, default=5, help="Number of clients to launch")
    parser.add_argument("--duration", type=int, default=60, help="Duration in seconds to run each client")
    parser.add_argument("--command", type=str, required=True, help="Command to launch the editor client")
    parser.add_argument("--max-concurrent", type=int, default=5,
                        help="Maximum number of clients to run concurrently (to avoid overwhelming the system)")
    args = parser.parse_args()
    
    print(f"Starting load test with {args.clients} clients for {args.duration} seconds each")
    print(f"Using command: {args.command}")
    print(f"Maximum concurrent clients: {args.max_concurrent}")
    
    start_time = time.time()
    
    # Set up a clean shutdown mechanism
    def signal_handler(sig, frame):
        print("\nInterrupted! Shutting down all clients...")
        sys.exit(0)
    
    signal.signal(signal.SIGINT, signal_handler)
    
    # Create a directory for logs
    os.chdir("analysis")
    os.makedirs("analysis_log", exist_ok=True)
    
    # Launch clients using a thread pool to manage concurrency
    with ThreadPoolExecutor(max_workers=args.max_concurrent) as executor:
        futures = [
            executor.submit(run_client, i, args.command, args.duration)
            for i in range(args.clients)
        ]
        
        # Wait for all clients to complete
        for future in futures:
            try:
                client_id = future.result()
            except Exception as e:
                print(f"Client execution failed: {e}")
    
    end_time = time.time()
    
    print(f"\nLoad test completed in {end_time - start_time:.2f} seconds")
    print(f"Client logs stored in {os.getcwd()}")

if __name__ == "__main__":
    main()