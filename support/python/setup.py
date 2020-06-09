import setuptools

setuptools.setup(
    name='nebula-sdk',
    use_scm_version={
        'root': '../..',
        'relative_to': __file__,
    },
    author='Puppet, Inc.',
    author_email='project-nebula-support@puppet.com',
    description='SDK for interacting with Project Nebula',
    url='https://github.com/puppetlabs/nebula-sdk',
    packages=setuptools.find_packages('nebula_sdk'),
    package_dir={'': 'nebula_sdk'},
    python_requires='>=3.8',
    setup_requires=[
        'setuptools_scm',
    ],
    install_requires=[
        'asgiref>=3.2.7',
        'hypercorn>=0.9.5',
        'requests>=2.23',
    ],
    classifiers=[
        'Intended Audience :: Developers',
        'License :: OSI Approved :: Apache Software License',
    ],
)
