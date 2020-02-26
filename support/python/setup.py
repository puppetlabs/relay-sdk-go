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
    packages=setuptools.find_packages('src'),
    package_dir={'': 'src'},
    python_requires='>=3.8',
    setup_requires=[
        'setuptools_scm',
    ],
    install_requires=[
        'requests>=2.23',
    ],
    classifiers=[
        'Intended Audience :: Developers',
        'License :: OSI Approved :: Apache Software License',
    ],
)
