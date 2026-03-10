from setuptools import setup, find_packages

setup(
    name="zero-lang",
    version="0.1.0",
    packages=find_packages(),
    install_requires=[],
    entry_points={
        "console_scripts": [
            "zero-compiler=zero_lang.compiler:main",
            "zero-generator=zero_lang.generator:main",
        ],
    },
    author="零点语言团队",
    description="零点语言编译器与代码生成器",
    long_description=open("README.md", encoding="utf-8").read(),
    long_description_content_type="text/markdown",
    url="https://github.com/zero-lang/zero-lang",
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires='>=3.6',
)