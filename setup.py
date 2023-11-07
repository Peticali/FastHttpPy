from setuptools import setup
import pathlib

HERE = pathlib.Path(__file__).parent
README = (HERE / "README.md").read_text()

setup(
    name='fasthttppy',
    version='0.0.2',
    packages=['fasthttppy'],
    package_data={"fasthttppy":["go_server/*.*"]},
    description="Python implementation of Golang FastHttp",
    url="https://github.com/Peticali/FastHttpPy",
    author="Pedro Palmeira",
    long_description_content_type="text/markdown",
    long_description=README,   
    author_email="pedropalmeira68@gmail.com",
    license="MIT",
)