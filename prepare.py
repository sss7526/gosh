import os
import yaml
import argparse
import logging
from pathlib import Path
from fnmatch import fnmatch

# Configure logging
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

# Mapping of file extensions to programming language names for Markdown code blocks
LANGUAGE_MAP = {
    '.py': 'python',
    '.js': 'javascript',
    '.ts': 'typescript',
    '.jsx': 'javascript',
    '.tsx': 'typescript',
    '.html': 'html',
    '.css': 'css',
    '.json': 'json',
    '.yaml': 'yaml',
    '.yml': 'yaml',
    '.java': 'java',
    '.c': 'c',
    '.cpp': 'cpp',
    '.cs': 'csharp',
    '.rb': 'ruby',
    '.go': 'go',
    '.php': 'php',
    '.sh': 'bash',
    '.bat': 'batch',
    '.sql': 'sql',
    '.swift': 'swift',
    '.kt': 'kotlin',
    '.rs': 'rust',
    '.xml': 'xml',
    '.ini': 'ini',
    '.md': 'markdown',
    '.txt': 'plaintext',
    '.svelte': 'svelte',
    '.tf': 'hcl',
    '.tfvars': 'hcl'
}

def load_config(config_path):
    """Load the configuration YAML file."""
    if not os.path.exists(config_path):
        logging.error(f"Configuration file '{config_path}' not found.")
        exit(1)

    try:
        with open(config_path, 'r') as config_file:
            return yaml.safe_load(config_file) or {}
    except yaml.YAMLError as e:
        logging.error(f"Error parsing YAML file: {e}")
        exit(1)
    except Exception as e:
        logging.error(f"Unexpected error while loading config: {e}")
        exit(1)

def format_code_block(code, language):
    """Format code as a Markdown code block."""
    return f"```{language}\n{code}\n```"

def detect_language(file_extension):
    """Detect the programming language based on file extension."""
    return LANGUAGE_MAP.get(file_extension, 'plaintext')

def process_file(file_path):
    """Read a file and format its content as a Markdown code block."""
    try:
        with open(file_path, 'r') as file:
            code = file.read()
        file_extension = Path(file_path).suffix
        language = detect_language(file_extension)
        return format_code_block(code, language)
    except Exception as e:
        logging.error(f"Error reading file '{file_path}': {e}")
        return None

def prepare_output_dir(output_dir):
    """Ensure the output directory exists."""
    try:
        Path(output_dir).mkdir(parents=True, exist_ok=True)
    except Exception as e:
        logging.error(f"Error creating output directory '{output_dir}': {e}")
        exit(1)

def should_include(file_path, include_patterns, exclude_patterns):
    """Check if a file should be included based on inclusion and exclusion patterns."""
    included = any(fnmatch(file_path, pattern) for pattern in include_patterns)
    excluded = any(fnmatch(file_path, pattern) for pattern in exclude_patterns)
    return included and not excluded

def generate_table_of_contents(sections):
    """Generate a Table of Contents for the Markdown document."""
    toc = ['# Table of Contents\n']
    for section in sections:
        anchor = section.lower().replace(' ', '-').replace('*', '').replace('#', '')
        toc.append(f"- [{section}](#{anchor})")
    return "\n".join(toc) + "\n\n"

def process_files(output, rules, output_dir):
    """Process files according to rules and generate Markdown."""
    prepare_output_dir(output_dir)

    output_path = os.path.join(output_dir, output)

    sections = []
    markdown_content = []

    for rule in rules:
        # Validate rule entries
        base_dir = rule.get('base_dir')
        if not base_dir or not os.path.exists(base_dir):
            logging.warning(f"Skipping invalid or missing base_dir: {base_dir}")
            continue

        include_patterns = rule.get('include', ['*'])
        exclude_patterns = rule.get('exclude', [])
        section_heading = rule.get('section_heading', None)

        # Add section heading if specified
        if section_heading:
            sections.append(section_heading)
            markdown_content.append(f"## {section_heading}\n\n")

        # Add description if specified
        description = rule.get('description', None)
        if description:
            markdown_content.append(f"> NOTE: {description}\n\n")

        # Process files in the directory
        for root, _, files in os.walk(base_dir):
            for file_name in files:
                file_path = os.path.join(root, file_name)
                relative_path = os.path.relpath(file_path, base_dir)

                if should_include(relative_path, include_patterns, exclude_patterns):
                    logging.info(f"Processing '{file_path}'...")
                    markdown_content.append(f"### File: `{relative_path}`\n")
                    block = process_file(file_path)
                    if block:
                        markdown_content.append(block + "\n\n")

    # Generate and add Table of Contents
    toc = generate_table_of_contents(sections)
    markdown_content.insert(0, toc)

    # Write the final Markdown file
    try:
        with open(output_path, 'w') as markdown_file:
            markdown_file.write('\n'.join(markdown_content))
        logging.info(f"Generated Markdown: '{output_path}'")
    except Exception as e:
        logging.error(f"Error writing to output file '{output_path}': {e}")
        exit(1)

def main():
    # Argument parsing
    parser = argparse.ArgumentParser(description="Generate Markdown documentation based on file processing rules.")
    parser.add_argument('--config', default='prepare_config.yaml', help="Path to the configuration YAML file (default: 'prepare_config.yaml').")
    parser.add_argument('--output_dir', default='output_docs', help="Directory to save the generated Markdown files (default: 'output_docs').")

    args = parser.parse_args()
    config_path = args.config
    output_dir = args.output_dir

    # Load the configuration
    config = load_config(config_path)

    # Validate configuration structure
    if 'outputs' not in config or not isinstance(config['outputs'], dict):
        logging.error("Invalid configuration format. 'outputs' section is missing or incorrect.")
        exit(1)

    # Process outputs
    for output, rules in config['outputs'].items():
        logging.info(f"Generating '{output}'...")
        if not isinstance(rules, list):
            logging.warning(f"Skipping output '{output}': Rules should be a list of dictionaries.")
            continue
        process_files(output, rules, output_dir)

    logging.info("Markdown documents prepared successfully.")

if __name__ == '__main__':
    main()