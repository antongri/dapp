module Dapp
  module Deployment
    module Dapp
      module Command
        module Apply
          def deployment_apply
            deployment.apps.each do |app|
              (app.kube.existing_deployments_names - app.to_kube_deployments.keys).each do |deployment_name|
                app.kube.delete_deployment!(deployment_name)
              end

              (app.kube.existing_services_names - app.to_kube_services.keys).each do |service_name|
                app.kube.delete_service!(service_name)
              end

              app.to_kube_deployments.each do |name, spec|
                if app.kube.deployment_exist?(name)
                  app.kube.replace_deployment!(name, spec) if app.kube.deployment_spec_changed?(name, spec)
                else
                  app.kube.create_deployment!(spec)
                end
              end

              app.to_kube_services.each do |name, spec|
                if app.kube.service_exist?(name)
                  app.kube.replace_service!(name, spec) if app.kube.service_spec_changed?(name, spec)
                else
                  app.kube.create_service!(spec)
                end
              end
            end
          end
        end
      end
    end
  end
end