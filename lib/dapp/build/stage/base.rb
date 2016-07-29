module Dapp
  module Build
    module Stage
      # Base of all stages
      class Base
        include Helper::Sha256
        include Helper::Trivia

        attr_accessor :prev_stage, :next_stage
        attr_reader :application

        def initialize(application, next_stage)
          @application = application

          @next_stage = next_stage
          @next_stage.prev_stage = self
        end

        # rubocop:disable Metrics/AbcSize, Metrics/MethodLength
        def build!
          return if should_be_skipped?
          prev_stage.build! if prev_stage
          begin
            if image.tagged?
              application.log_state(name, state: application.t(code: 'state.using_cache'))
            elsif application.dry_run?
              application.log_state(name, state: application.t(code: 'state.build'), styles: { status: :success })
            else
              application.log_process(name, process: application.t(code: 'status.process.building'), short: should_be_not_detailed?) do
                image_build!
              end
            end
          ensure
            log_build
          end
          raise Exception::IntrospectImage,
                message: application.t(code: 'introspect.stage', data: { name: name }),
                data: { built_id: image.built_id, options: image.send(:prepared_options) } if should_be_introspected?
        end
        # rubocop:enable Metrics/AbcSize, Metrics/MethodLength

        def save_in_cache!
          return if image.tagged?
          prev_stage.save_in_cache!                                                          if prev_stage
          image.tag!(log_verbose: application.log_verbose?, log_time: application.log_time?) unless application.dry_run?
        end

        def signature
          hashsum prev_stage.signature
        end

        def image
          @image ||= begin
            Image::Stage.new(name: image_name, from: from_image).tap do |image|
              image.add_volume "#{application.tmp_path}:#{application.container_tmp_path}"
              yield image if block_given?
            end
          end
        end

        protected

        def name
          class_to_lowercase.to_sym
        end

        def should_be_not_detailed?
          image.send(:bash_commands).empty?
        end

        def should_be_skipped?
          image.tagged? && !application.log_verbose? && application.cli_options[:introspect_stage].nil?
        end

        def should_be_introspected?
          application.cli_options[:introspect_stage] == name && !application.dry_run?
        end

        def image_build!
          image.build!(log_verbose: application.log_verbose?,
                       log_time: application.log_time?,
                       introspect_error: application.cli_options[:introspect_error],
                       introspect_before_error: application.cli_options[:introspect_before_error])
        end

        def from_image
          prev_stage.image if prev_stage || begin
            raise Error::Build, code: :from_image_required
          end
        end

        def image_name
          "dapp:#{signature}"
        end

        def image_info
          date, size = image.info
          _date, from_size = from_image.info
          [date, (from_size.to_f - size.to_f).abs]
        end

        def format_image_info
          date, size = image_info
          application.t(code: 'image.info', data: { date: Time.parse(date).localtime, size: size.to_f.round(2) })
        end

        # rubocop:disable Metrics/AbcSize
        def log_build
          application.with_log_indent do
            application.log_info application.t(code: 'image.signature', data: { signature: image_name })
            application.log_info format_image_info if image.tagged?
            unless (bash_commands = image.send(:bash_commands)).empty?
              application.log_info application.t(code: 'image.commands')
              application.with_log_indent { application.log_info bash_commands.join("\n") }
            end
          end if application.log? && application.log_verbose?
        end
      end # Base
    end # Stage
  end # Build
end # Dapp
