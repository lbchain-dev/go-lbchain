Pod::Spec.new do |spec|
  spec.name         = 'Glbchain-dev'
  spec.version      = '{{.Version}}'
  spec.license      = { :type => 'GNU Lesser General Public License, Version 3.0' }
  spec.homepage     = 'https://github.com/lbchain-devchain/go-lbchain-dev'
  spec.authors      = { {{range .Contributors}}
		'{{.Name}}' => '{{.Email}}',{{end}}
	}
  spec.summary      = 'iOS lbchain-devchain Client'
  spec.source       = { :git => 'https://github.com/lbchain-devchain/go-lbchain-dev.git', :commit => '{{.Commit}}' }

	spec.platform = :ios
  spec.ios.deployment_target  = '9.0'
	spec.ios.vendored_frameworks = 'Frameworks/Glbchain-dev.framework'

	spec.prepare_command = <<-CMD
    curl https://glbchain-devstore.blob.core.windows.net/builds/{{.Archive}}.tar.gz | tar -xvz
    mkdir Frameworks
    mv {{.Archive}}/Glbchain-dev.framework Frameworks
    rm -rf {{.Archive}}
  CMD
end
