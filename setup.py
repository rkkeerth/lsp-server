#!/usr/bin/env python3
"""Setup script for the basic LSP server."""

from setuptools import setup, find_packages

with open("README.md", "r", encoding="utf-8") as fh:
    long_description = fh.read()

setup(
    name="basic-lsp-server",
    version="0.1.0",
    author="LSP Server Developer",
    description="A basic Language Server Protocol implementation using only Python standard library",
    long_description=long_description,
    long_description_content_type="text/markdown",
    py_modules=["server"],
    python_requires=">=3.8",
    install_requires=[
        # No external dependencies - uses only Python standard library
    ],
    entry_points={
        "console_scripts": [
            "basic-lsp-server=server:main",
        ],
    },
    classifiers=[
        "Development Status :: 3 - Alpha",
        "Intended Audience :: Developers",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "Programming Language :: Python :: 3.12",
    ],
)
