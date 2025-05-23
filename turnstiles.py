#!/usr/bin/env python3
"""
Raspberry Pi Keyboard Input Validator with GPIO Relay Control

This script listens for keyboard input, validates codes via API call,
and triggers a relay when valid codes are entered.

Requirements:
- gpiod library for Raspberry Pi GPIO control
- requests library for API calls
- Python 3.x

Hardware Setup:
- Connect relay to GPIO pin (default: GPIO 18)
- Ensure proper relay wiring with appropriate power supply
"""

import gpiod
import time
import sys
import requests
import threading
import configparser
import os
from typing import Optional

# Configuration
CONFIG_FILE = "config.ini"  # Configuration file path
DEFAULT_API_ENDPOINT = "https://beta-backend-dev-kpe3ohblca-ew.a.run.app/v2/turnstiles/doConsume"  # Replace with your API endpoint
DEFAULT_RELAY_PIN = 17  # GPIO pin for relay control

API_TIMEOUT = 5.0  # API request timeout in seconds
RELAY_ACTIVE_TIME = 0.2  # How long to keep relay active (seconds)



class InputValidator:
    def __init__(self, relay_pin: int = DEFAULT_RELAY_PIN, api_url: str = DEFAULT_API_ENDPOINT, config_file: str = CONFIG_FILE):
        self.relay_pin = relay_pin
        self.api_url = api_url
        self.config_file = config_file
        self.api_key = None
        self.input_buffer = ""
        self.load_config()
        self.setup_gpio()

    def load_config(self):
        """Load configuration from config file"""
        try:
            if not os.path.exists(self.config_file):
                print(f"Config file '{self.config_file}' not found. Creating default config...")
                self.create_default_config()
                return

            config = configparser.ConfigParser()
            config.read(self.config_file)

            # Load API settings
            if 'API' in config:
                self.api_key = config['API'].get('api_key', '').strip()

                # Override API URL if specified in config
                config_api_url = config['API'].get('endpoint', '').strip()
                if config_api_url:
                    self.api_url = config_api_url

                # Override timeout if specified
                try:
                    timeout = config['API'].getfloat('timeout')
                    if timeout > 0:
                        global API_TIMEOUT
                        API_TIMEOUT = timeout
                except (ValueError, TypeError):
                    pass

            # Load GPIO settings
            if 'GPIO' in config:
                try:
                    config_relay_pin = config['GPIO'].getint('relay_pin')
                    if config_relay_pin > 0:
                        self.relay_pin = config_relay_pin
                except (ValueError, TypeError):
                    pass

                try:
                    active_time = config['GPIO'].getfloat('relay_active_time')
                    if active_time > 0:
                        global RELAY_ACTIVE_TIME
                        RELAY_ACTIVE_TIME = active_time
                except (ValueError, TypeError):
                    pass

            if self.api_key:
                print("✓ Configuration loaded successfully")
            else:
                print("⚠ Warning: No API key found in config file")

        except Exception as e:
            print(f"Error loading config: {e}")
            print("Using default settings...")

    def create_default_config(self):
        """Create a default configuration file"""
        config = configparser.ConfigParser()

        # API Configuration
        config['API'] = {
            'api_key': 'YOUR_API_KEY_HERE',
            'endpoint': DEFAULT_API_ENDPOINT,
            'timeout': '5.0'
        }

        # GPIO Configuration
        config['GPIO'] = {
            'relay_pin': DEFAULT_RELAY_PIN,
            'relay_active_time': '0.2'
        }

        try:
            with open(self.config_file, 'w') as f:
                config.write(f)
            print(f"✓ Default config file created: {self.config_file}")
            print("Please edit the config file and add your API key before running again.")
        except Exception as e:
            print(f"Error creating config file: {e}")

    def setup_gpio(self):
        """Initialize GPIO settings for relay control"""
        try:
            chip = gpiod.Chip('gpiochip4')
            self.relay_line = chip.get_line(self.relay_pin)
            self.relay_line.request(consumer="LED", type=gpiod.LINE_REQ_DIR_OUT)
            print(f"GPIO initialized - Relay pin: {self.relay_pin}")
        except Exception as e:
            print(f"GPIO setup error: {e}")

    def validate_code_api(self, code: str) -> bool:
        """
        Validate code via API call

        Args:
            code: The input code to validate

        Returns:
            bool: True if code is valid, False otherwise
        """
        try:
            # Prepare API request
            payload = {
                "guid": code+":0", # 1 for dir out todo make optional
                "allowReentry": True, # todo make configurable
            }

            headers = {
                "Content-Type": "application/json",
                "Authorization": self.api_key
            }

            print(f"Validating code: {code}")

            # Make API request
            response = requests.post(
                self.api_url,
                json=payload,
                headers=headers,
                timeout=API_TIMEOUT
            )
            # Check response
            if response.status_code == 200:
                result = response.json()
                print(f"response {result}")
                e = result.get("error")
#                 is_valid = result.get("valid", False)
#                 print(f"API Response: {'Valid' if is_valid else 'Invalid'}")
#                 return is_valid
                if (e is None or e.get("message") == "success"):
                    return True
                else:
                    return False
            else:
                result = response.json()
                print(f"response {result}")
                print(f"API Error: Status {response.status_code}")
                return False

        except requests.exceptions.Timeout:
            print("API request timeout")
            return False
        except requests.exceptions.RequestException as e:
            print(f"API request error: {e}")
            return False
        except Exception as e:
            print(f"Validation error: {e}")
            return False

    def trigger_relay(self):
        """Activate relay for specified duration"""
        try:
            print("Activating relay...")
            self.relay_line.set_value(1)
            time.sleep(RELAY_ACTIVE_TIME)
            print("Deactivating relay...")
            self.relay_line.set_value(0)

        except Exception as e:
            print(f"Relay control error: {e}")
            # Ensure relay is turned off in case of error
            try:
                self.relay_line.set_value(0)
            except:
                pass

    def process_input(self, code: str):
        """Process the input code through validation and relay control"""
        if not code.strip():
            return

        print(f"\nProcessing input: '{code}'")

        # Validate code via API
        if self.validate_code_api(code):
            print("✓ Code validated successfully!")

            # Trigger relay in separate thread to avoid blocking
            relay_thread = threading.Thread(target=self.trigger_relay)
            relay_thread.daemon = True
            relay_thread.start()

        else:
            print("✗ Invalid code")

    def listen_for_input(self):
        """Main input listening loop"""
        print("Keyboard Input Validator Started")
        print("Enter codes and press ENTER to validate")
        print("Type 'quit' or 'exit' to stop")
        print("-" * 40)

        try:
            while True:
                try:
                    # Get input from user
                    user_input = input("Enter code: ").strip()

                    # Check for exit commands
                    if user_input.lower() in ['quit', 'exit', 'q']:
                        print("Exiting...")
                        break
                    if user_input.lower() in ['test']:
                        print("Testing relay...")
                        self.trigger_relay()
                    # Process the input
                    elif user_input:
                        self.process_input(user_input)

                except KeyboardInterrupt:
                    print("\nKeyboard interrupt received. Exiting...")
                    break
                except EOFError:
                    print("\nEOF received. Exiting...")
                    break

        finally:
            self.cleanup()

    def cleanup(self):
        """Clean up GPIO and other resources"""
        try:
            self.relay_line.release()
            print("GPIO cleaned up")
        except Exception as e:
            print(f"GPIO cleanup error: {e}")

def main():
    """Main function to run the input validator"""

    # Create and start the validator
    validator = InputValidator()
    validator.listen_for_input()

if __name__ == "__main__":
    main()